import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {LoginComponent} from './login.component';
import {FourOFourComponent} from '../four-ofour/four-ofour.component';

export const spellbookRoutes: Routes = [
	{
		path: '',
		redirectTo: 'content',
		pathMatch: 'full'
	},
	{
		path: 'login',
		component: LoginComponent
	},
	{
		path: 'fourofour',
		component: FourOFourComponent
	}
];

@NgModule({
	imports: [
		RouterModule.forChild(spellbookRoutes),
	],
	exports: [RouterModule]
})
export class SpellbookRoutingModule {
}
