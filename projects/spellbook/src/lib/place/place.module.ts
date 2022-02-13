import {ModuleWithProviders, NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ListPlaceComponent} from './list-place.component';
import {CreateUpdatePlaceComponent} from './create-update-place.component';
import {FlexLayoutModule} from '@angular/flex-layout';
import {RouterModule, Routes} from '@angular/router';
import {ReactiveFormsModule} from '@angular/forms';
import {
	MatAutocompleteModule,
	MatButtonModule,
	MatCardModule,
	MatCheckboxModule,
	MatChipsModule,
	MatDialogModule,
	MatExpansionModule,
	MatFormFieldModule,
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
import {SnackbarComponent} from '../core/snackbar.component';
import {AgmCoreModule, GoogleMapsAPIWrapper} from '@agm/core';
import {AuthService} from '../core/auth.service';
import {UtilsModule} from '../core/utils.module';

@NgModule({
	declarations: [
		ListPlaceComponent,
		CreateUpdatePlaceComponent
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
		RouterModule.forChild(PlaceModule.routes),
		AgmCoreModule,
		UtilsModule
	],
	entryComponents: [
		SnackbarComponent
	],
})
export class PlaceModule {

	static readonly routes: Routes = [
		{
			path: '',
			component: ListPlaceComponent,
			canActivate: [AuthService]
		},
		{
			path: ':id',
			component: CreateUpdatePlaceComponent,
			canActivate: [AuthService]
		}
	];

	static forRoot(token: string): ModuleWithProviders {
		const providers = AgmCoreModule.forRoot({
				apiKey: token,
				libraries: ['places']
			}
		).providers;
		return {
			ngModule: PlaceModule,
			providers: [
				providers,
				GoogleMapsAPIWrapper
			]
		};
	}
}
