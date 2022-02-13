import {NgModule} from '@angular/core';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {AppRoutingModule} from './app-routing.module';
import {TestModule} from './test/test.module';
import {TestComponent} from './test.component';
import {
	MatButtonModule,
	MatIconModule,
	MatListModule,
	MatMenuModule,
	MatSidenavModule,
	MatToolbarModule
} from '@angular/material';
import {SpellbookModule} from '../../../spellbook/src/lib/core/spellbook.module';
import {Definition} from '../../../spellbook/src/lib/core/typedef.service';
import {EventDefinition, SpecialDefinition} from './defs/content-definitions';
import {DragulaModule} from 'ng2-dragula';

export function typeDefinitions(): Array<Definition<any>> {
	return [
		new EventDefinition(),
		new SpecialDefinition()
	];
}

@NgModule({
	declarations: [
		AppComponent,
		TestComponent,
	],
	imports: [
		BrowserAnimationsModule,
		AppRoutingModule,
		SpellbookModule.forRoot({
			apiUrl: 'http://localhost:4200/api',
			superUserRedirectUrl: '/api/superuser/',
			googleMapKey: 'AIzaSyCoeugGhGD6Gxwo1VtyQ9G3dZGnNY9bfQs',
			typeDefinitions,
		}),
		DragulaModule.forRoot(),
		TestModule,
		MatButtonModule,
		MatToolbarModule,
		MatSidenavModule,
		MatListModule,
		MatIconModule,
		MatMenuModule,
	],
	providers: [],
	bootstrap: [AppComponent]
})
export class AppModule {
}
