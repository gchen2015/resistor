//Auth examples from: https://github.com/auth0-blog/angular2-authentication-sample
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule } from '@angular/router';
import { HttpModule } from '@angular/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { DropdownModule } from 'ng2-bootstrap';

// external libs

import { Ng2TableModule } from './common/ng-table/ng2-table';

import { Home } from './home/home';
import { Login } from './login/login';
import { App } from './app/app';
import { HttpAPI } from './common/httpAPI';

import { AppRoutes } from './app/app.routes';
//common
import { ControlMessagesComponent } from './common/control-messages.component';
import { MultiselectDropdownModule } from './common/multiselect-dropdown'
import { PasswordToggleDirective } from './common/custom-directives'
import { TableActions } from './common/table-actions';

//Resistor Components

import { BlockUIService } from './common/blockui/blockui-service';

import { AccordionModule , PaginationModule ,TabsModule } from 'ng2-bootstrap';
import { TooltipModule } from 'ng2-bootstrap';
import { ModalModule } from 'ng2-bootstrap';
import { ModalDirective } from 'ng2-bootstrap';
import { ProgressbarModule } from 'ng2-bootstrap';

import { GenericModal } from './common/generic-modal';
import { AboutModal } from './home/about-modal';
import { TreeView } from './common/dataservice/treeview';
import { ImportFileModal } from './common/dataservice/import-file-modal'

//others
import { WindowRef } from './common/windowref';
import { ValidationService } from './common/validation.service';
import { ExportServiceCfg } from './common/dataservice/export.service'
//pipes
import { ObjectParserPipe,SplitCommaPipe } from './common/custom_pipe';
import { ElapsedSecondsPipe } from './common/elapsedseconds.pipe';

import { CustomPipesModule } from './common/custom-pipe-module';

import { BlockUIComponent } from './common/blockui/blockui-component';
import { SpinnerComponent } from './common/spinner';


@NgModule({
  bootstrap: [App],
  declarations: [
    PasswordToggleDirective,
    ObjectParserPipe,
    //ElapsedSecondsPipe,
    SplitCommaPipe,
    TableActions,
    ControlMessagesComponent,
    GenericModal,
    AboutModal,
    ImportFileModal,
    BlockUIComponent,
    TreeView,
    SpinnerComponent,
    Home,
    Login,
    App,
  ],
  imports: [
    CustomPipesModule,
    HttpModule,
    BrowserModule,
    FormsModule,
    ReactiveFormsModule,
    MultiselectDropdownModule,
    ProgressbarModule.forRoot(),
    AccordionModule.forRoot(),
    TooltipModule.forRoot(),
    ModalModule.forRoot(),
    PaginationModule.forRoot(),
    TabsModule.forRoot(),
    DropdownModule.forRoot(),
    Ng2TableModule,
    RouterModule.forRoot(AppRoutes)
  ],
  providers: [
    WindowRef,
    HttpAPI,
    ExportServiceCfg,
    ValidationService,
    BlockUIService
  ],
  entryComponents: [BlockUIComponent]
})
export class AppModule {}
