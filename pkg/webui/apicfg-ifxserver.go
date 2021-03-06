package webui

import (
	"fmt"
	"time"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"

	"github.com/influxdata/influxdb/client/v2"
)

// NewAPICfgIfxServer API for IfxServer Catalog Management
func NewAPICfgIfxServer(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/ifxserver", func() {
		m.Get("/", reqSignedIn, GetIfxServer)
		m.Post("/", reqSignedIn, bind(config.IfxServerCfg{}), AddIfxServer)
		m.Put("/:id", reqSignedIn, bind(config.IfxServerCfg{}), UpdateIfxServer)
		m.Delete("/:id", reqSignedIn, DeleteIfxServer)
		m.Get("/:id", reqSignedIn, GetIfxServerCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetIfxServerAffectOnDel)
		m.Post("/ping", reqSignedIn, bind(config.IfxServerCfg{}), PingIfxServer)
		m.Post("/import", reqSignedIn, bind(config.IfxServerCfg{}), ImportIfxCatalog)
	})

	return nil
}

// ping internal
func ping(cfg *config.IfxServerCfg) (client.Client, time.Duration, string, error) {

	conf := client.HTTPConfig{
		Addr:      cfg.URL,
		Username:  cfg.AdminUser,
		Password:  cfg.AdminPasswd,
		UserAgent: "resistor",
		//		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}
	cli, err := client.NewHTTPClient(conf)

	if err != nil {
		return cli, 0, "", err
	}
	elapsed, message, err := cli.Ping(time.Duration(5) * time.Second)
	if err == nil {
		_, err = InfluxQuery(cli, "show databases", "")
	}
	return cli, elapsed, message, err
}

// InfluxQuery does IQL queries
func InfluxQuery(cli client.Client, query string, db string) ([]string, error) {
	defer cli.Close()
	q := client.NewQuery(query, db, "s")
	response, err := cli.Query(q)
	if err != nil {
		log.Errorf("Query RESULT ERROR: %+s", err)
		return nil, err
	}
	if response.Error() != nil {
		log.Errorf("Query RESULT ERROR: %+s", response.Error())
		return nil, fmt.Errorf("Query RESULT ERROR: %+s", response.Error())
	}
	var retarray []string
	for _, v1 := range response.Results {
		for _, v2 := range v1.Series {
			for _, v3 := range v2.Values {
				retarray = append(retarray, v3[0].(string))
			}
		}
	}
	return retarray, nil
}

func getInfluxDBs(cli client.Client) []string {
	ret, err := InfluxQuery(cli, "SHOW DATABASES", "")
	if err != nil {
		log.Errorf("GetInfluxDB's Query: %s", err)
	}
	return ret
}

func getDBMeasurements(cli client.Client, db string) []string {
	ret, err := InfluxQuery(cli, "SHOW MEASUREMENTS", db)
	if err != nil {
		log.Errorf("Get Measurements Query: %s", err)
	}
	return ret
}

func getDBRetention(cli client.Client, db string) []string {
	ret, err := InfluxQuery(cli, "SHOW RETENTION POLICIES ON \""+db+"\"", "")
	if err != nil {
		log.Errorf("Get Retention Policies Query: %s", err)
	}
	return ret
}

func getMeasurementsFields(cli client.Client, db string, m string) []string {

	ret, err := InfluxQuery(cli, "SHOW FIELD KEYS FROM \""+m+"\"", db)
	if err != nil {
		log.Errorf("Get Measurements Fields Query: %s", err)
	}
	return ret
}

func getMeasurementsTags(cli client.Client, db string, m string) []string {
	ret, err := InfluxQuery(cli, "SHOW TAG KEYS FROM \""+m+"\"", db)
	if err != nil {
		log.Errorf("Get Measurements Tags Query: %s", err)
	}
	return ret
}

// ImportIfxCatalog new snmpdevice to de internal BBDD --pending--
func ImportIfxCatalog(ctx *Context, dev config.IfxServerCfg) {
	log.Warningf("Importing catalog for Server: %s", dev.ID)
	cli, _, _, err := ping(&dev)
	if err != nil {
		log.Errorf("Error on Ping Influx Server %s: Err: %s", dev.ID, err)
		ctx.JSON(404, err.Error())
		return
	}

	dbs := getInfluxDBs(cli)
	for _, db := range dbs {
		rps := getDBRetention(cli, db)
		meas := getDBMeasurements(cli, db)

		var itemarray []*config.ItemComponent
		for _, m := range meas {

			tags := getMeasurementsTags(cli, db, m)
			fields := getMeasurementsFields(cli, db, m)
			mcfg := config.IfxMeasurementCfg{Name: m, Tags: tags, Fields: fields}
			_, err := agent.MainConfig.Database.AddIfxMeasurementCfg(&mcfg)
			if err != nil {
				log.Errorf("Error on Importing Influx DBs: %s Err: %s", dev.ID, err)
				ctx.JSON(404, err.Error())
				return
			}
			log.Infof("Got DATABASE [%s] with retention policy [%s] MEASUREMENT %d/%s TAGS [%+v] FIELDS [%+v]", db, rps, mcfg.ID, m, tags, fields)
			itemarray = append(itemarray, &config.ItemComponent{ID: mcfg.ID, Name: m})
		}

		dbcfg := config.IfxDBCfg{Name: db, IfxServer: dev.ID, Retention: rps, Measurements: itemarray}
		_, err := agent.MainConfig.Database.AddOrUpdateIfxDBCfg(&dbcfg)
		if err != nil {
			log.Errorf("Error on Importing Influx DBs: %s Err: %s", dev.ID, err)
			ctx.JSON(404, err.Error())
			return
		}
	}

	ctx.JSON(200, &dbs)
}

// PingIfxServer Insert new snmpdevice to de internal BBDD --pending--
func PingIfxServer(ctx *Context, dev config.IfxServerCfg) {
	_, elapsed, message, err := ping(&dev)
	if err != nil {
		log.Warningf("Error on Ping Influx Server %s: Err: %s", dev.ID, err)
		ctx.JSON(404, err.Error())
		return
	}
	ctx.JSON(200, &struct {
		Message string
		Elapsed time.Duration
	}{
		Message: message,
		Elapsed: elapsed,
	})

}

// GetIfxServer Return snmpdevice list to frontend
func GetIfxServer(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetIfxServerCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddIfxServer Insert new snmpdevice to de internal BBDD --pending--
func AddIfxServer(ctx *Context, dev config.IfxServerCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddIfxServerCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateIfxServer --pending--
func UpdateIfxServer(ctx *Context, dev config.IfxServerCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateIfxServerCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteIfxServer delete from the catalog database
func DeleteIfxServer(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelIfxServerCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetIfxServerCfgByID --pending--
func GetIfxServerCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetIfxServerCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetIfxServerAffectOnDel --pending--
func GetIfxServerAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetIfxServerCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
