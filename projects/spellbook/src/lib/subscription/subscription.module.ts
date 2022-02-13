import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ListSubscriptionComponent} from './list-subscription.component';
import {SnackbarComponent} from '../core/snackbar.component';
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
import {AuthService} from '../core/auth.service';
import {CreateUpdateSubscriptionComponent} from './create-update-subscription.component';

@NgModule({
	declarations: [ListSubscriptionComponent, CreateUpdateSubscriptionComponent],
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
		RouterModule.forChild(SubscriptionModule.routes),
	],
	entryComponents: [
		SnackbarComponent
	],
})
export class SubscriptionModule {

	static readonly routes: Routes = [
		{
			path: '',
			component: ListSubscriptionComponent,
			canActivate: [AuthService]
		},
		{
			path: ':key',
			component: CreateUpdateSubscriptionComponent,
			canActivate: [AuthService]
		}
	];
}
