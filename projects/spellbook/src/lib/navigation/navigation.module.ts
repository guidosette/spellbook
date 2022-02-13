import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ListPageComponent} from './list-page.component';
import {CreateUpdatePageComponent} from './create-update-page.component';
import {
	MatButtonModule,
	MatCardModule,
	MatFormFieldModule,
	MatIconModule,
	MatInputModule,
	MatPaginatorModule,
	MatProgressSpinnerModule,
	MatSelectModule,
	MatSortModule,
	MatTableModule
} from '@angular/material';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {FlexLayoutModule} from '@angular/flex-layout';
import {RouterModule, Routes} from '@angular/router';
import {AuthService} from '../core/auth.service';

@NgModule({
	declarations: [
		ListPageComponent,
		CreateUpdatePageComponent
	],
	imports: [
		CommonModule,
		MatIconModule,
		MatProgressSpinnerModule,
		MatPaginatorModule,
		MatFormFieldModule,
		MatTableModule,
		MatSortModule,
		MatCardModule,
		FormsModule,
		ReactiveFormsModule,
		MatButtonModule,
		FlexLayoutModule,
		MatInputModule,
		MatSelectModule,
		RouterModule.forChild(NavigationModule.routes)
	],
	entryComponents: []
})
export class NavigationModule {

	static readonly routes: Routes = [
		{
			path: '',
			component: ListPageComponent,
			canActivate: [AuthService]
		},
		{
			path: ':id',
			component: CreateUpdatePageComponent,
			canActivate: [AuthService]
		}
	];
}
