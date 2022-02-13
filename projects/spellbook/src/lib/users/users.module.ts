import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {FlexLayoutModule} from '@angular/flex-layout';
import {CommonModule} from '@angular/common';
import {
	MatButtonModule,
	MatCardModule,
	MatCheckboxModule,
	MatExpansionModule,
	MatFormFieldModule,
	MatIconModule,
	MatInputModule,
	MatPaginatorModule,
	MatProgressSpinnerModule,
	MatSelectModule,
	MatSidenavModule,
	MatSortModule,
	MatTableModule,
	MatToolbarModule
} from '@angular/material';
import {ListUserComponent} from './list-user.component';
import {CreateUpdateUserComponent} from './create-update-user.component';
import {ReactiveFormsModule} from '@angular/forms';
import {AuthService} from '../core/auth.service';

@NgModule({
	declarations: [
		CreateUpdateUserComponent,
		ListUserComponent,
	],
	imports: [
		CommonModule,
		FlexLayoutModule,
		RouterModule,
		MatSidenavModule,
		MatSortModule,
		MatProgressSpinnerModule,
		MatToolbarModule,
		MatTableModule,
		MatCardModule,
		MatFormFieldModule,
		MatButtonModule,
		MatPaginatorModule,
		MatIconModule,
		MatInputModule,
		MatCheckboxModule,
		MatSelectModule,
		MatExpansionModule,
		ReactiveFormsModule,
		RouterModule.forChild(UsersModule.routes)
	],
	providers: []
})

export class UsersModule {

	static readonly routes: Routes = [
		{
			path: '',
			component: ListUserComponent,
			canActivate: [AuthService]
		},
		{
			path: ':username',
			component: CreateUpdateUserComponent,
			canActivate: [AuthService]
		}
	];
}
