import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CKEditorModule} from '@ckeditor/ckeditor5-angular';
import {
	MAT_DATE_LOCALE,
	MatAutocompleteModule,
	MatButtonModule,
	MatCardModule,
	MatCheckboxModule,
	MatChipsModule,
	MatDialogModule,
	MatExpansionModule,
	MatFormFieldModule,
	MatGridListModule,
	MatIconModule,
	MatInputModule,
	MatMenuModule,
	MatPaginatorModule,
	MatProgressSpinnerModule,
	MatRadioModule,
	MatSelectModule,
	MatSidenavModule,
	MatSnackBarModule,
	MatSortModule,
	MatTableModule,
	MatToolbarModule
} from '@angular/material';
import {ReactiveFormsModule} from '@angular/forms';
import {CreateUpdateContentComponent} from './create-update-content.component';
import {ListContentComponent} from './list-content.component';
import {FlexLayoutModule} from '@angular/flex-layout';
import {RouterModule, Routes} from '@angular/router';
import {MediaModule} from '../media/media.module';


import {OWL_DATE_TIME_LOCALE, OwlDateTimeModule, OwlNativeDateTimeModule} from 'ng-pick-datetime';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {BrowserModule} from '@angular/platform-browser';
import {AttachmentComponent} from '../media/multimedia/attachment.component';
import {CreateAttachmentGroupComponent} from '../media/multimedia/create-attachment-group.component';
import {SnackbarComponent} from '../core/snackbar.component';
import {AuthService} from '../core/auth.service';
import {UtilsModule} from '../core/utils.module';

@NgModule({
	declarations: [
		CreateUpdateContentComponent,
		ListContentComponent
	],
	imports: [
		CKEditorModule,
		CommonModule,
		FlexLayoutModule,
		RouterModule,
		ReactiveFormsModule,
		MatSidenavModule,
		MatSortModule,
		MatProgressSpinnerModule,
		MatToolbarModule,
		MatTableModule,
		MatCardModule,
		MatFormFieldModule,
		MatButtonModule,
		MatPaginatorModule,
		MatInputModule,
		MatCheckboxModule,
		MatAutocompleteModule,
		MatSelectModule,
		MatChipsModule,
		MatIconModule,
		MatDialogModule,
		MatMenuModule,
		MatExpansionModule,
		MatRadioModule,
		MatSnackBarModule,
		MatGridListModule,
		OwlDateTimeModule,
		OwlNativeDateTimeModule,
		BrowserModule,
		BrowserAnimationsModule,
		RouterModule.forChild(ContentModule.routes),
		UtilsModule,
		MediaModule,
	],
	entryComponents: [
		AttachmentComponent,
		CreateAttachmentGroupComponent,
		SnackbarComponent
	],
	providers: [
		{provide: MAT_DATE_LOCALE, useValue: 'it-IT'},
		{provide: OWL_DATE_TIME_LOCALE, useValue: 'it'},
	],
	exports: []
})
export class ContentModule {

	static readonly routes: Routes = [
		{
			path: '',
			component: ListContentComponent,
			canActivate: [AuthService]
		},
		{
			path: ':id',
			component: CreateUpdateContentComponent,
			canActivate: [AuthService]
		}
	];
}
