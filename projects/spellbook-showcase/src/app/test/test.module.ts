import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {TestListComponent} from './test-list/test-list.component';
import {TestDetailComponent} from './test-detail/test-detail.component';
import {RouterModule, Routes} from '@angular/router';
import {AuthService} from '../../../../spellbook/src/lib/core/auth.service';

@NgModule({
	declarations: [TestListComponent, TestDetailComponent],
	imports: [
		CommonModule,
		RouterModule.forChild(TestModule.routes),

	]
})
export class TestModule {

	static readonly routes: Routes = [
		{
			path: '',
			component: TestListComponent,
			canActivate: [AuthService]
		},
		{
			path: ':id',
			component: TestDetailComponent,
			canActivate: [AuthService]
		}
	];
}
