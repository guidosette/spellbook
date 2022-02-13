import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ListMediaComponent} from './list-media.component';
import {FlexLayoutModule} from '@angular/flex-layout';
import {RouterModule, Routes} from '@angular/router';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {
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
import {CreateUpdateMediaComponent} from './create-update-media.component';
import {BrowseMediaFilesComponent} from './browse-media-files.component';
import {ShowFullscreenMediaComponent} from './show-fullscreen-media.component';
import {UploadMediaFileComponent} from './upload-media-file.component';
import {CollectMediaUrlComponent} from './collect-media-url.component';
import {ListMediaDialogComponent} from './list-media-dialog.component';
import {AuthService} from '../core/auth.service';
import {AttachmentGroupComponent} from './multimedia/attachment-group.component';
import {CreateAttachmentGroupComponent} from './multimedia/create-attachment-group.component';
import {AttachmentComponent} from './multimedia/attachment.component';
import {DragDropModule} from '@angular/cdk/drag-drop';
import {DragulaModule} from 'ng2-dragula';
import {MatTooltipModule} from '@angular/material/tooltip';

@NgModule({
	declarations: [
		ListMediaComponent,
		ListMediaDialogComponent,
		CreateUpdateMediaComponent,
		BrowseMediaFilesComponent,
		ShowFullscreenMediaComponent,
		UploadMediaFileComponent,
		CollectMediaUrlComponent,
		AttachmentComponent,
		AttachmentGroupComponent,
		CreateAttachmentGroupComponent
	],
	imports: [
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
		FormsModule,
		RouterModule.forChild(MediaModule.routes),
		DragDropModule,
		DragulaModule.forRoot(),
		MatTooltipModule,
	],
	exports: [
		AttachmentGroupComponent,
		AttachmentComponent,
		CreateAttachmentGroupComponent,
	],
	entryComponents: [
		CreateUpdateMediaComponent,
		BrowseMediaFilesComponent,
		ShowFullscreenMediaComponent,
		UploadMediaFileComponent,
		CollectMediaUrlComponent,
		ListMediaComponent,
		ListMediaDialogComponent,
		AttachmentComponent,
		CreateAttachmentGroupComponent,
		AttachmentGroupComponent,
	]
})
export class MediaModule {

	static readonly routes: Routes = [
		{
			path: '',
			component: ListMediaComponent,
			canActivate: [AuthService]
		},
		{
			path: ':id',
			component: ShowFullscreenMediaComponent,
			canActivate: [AuthService]
		}
	];
}
