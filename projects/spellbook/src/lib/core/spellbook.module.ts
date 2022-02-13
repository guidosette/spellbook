import {APP_INITIALIZER, ModuleWithProviders, NgModule} from '@angular/core';
import {LoginComponent} from './login.component';
import {Spellbook} from './spellbook';
import {SpellbookConfig} from './spellbook-config';
import {CommonModule} from '@angular/common';
import {HttpClientModule} from '@angular/common/http';
import {
	MAT_SNACK_BAR_DEFAULT_OPTIONS,
	MatButtonModule,
	MatCardModule,
	MatFormFieldModule,
	MatIconModule,
	MatInputModule,
	MatListModule,
	MatMenuModule,
	MatSelectModule,
	MatSidenavModule,
	MatToolbarModule
} from '@angular/material';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {FlexLayoutModule} from '@angular/flex-layout';
import {SpellbookRoutingModule} from './spellbook-routing.module';
import {SpellbookComponent} from './spellbook.component';
import {UsersModule} from '../users/users.module';
import {AuthService} from './auth.service';
import {FourOFourComponent} from '../four-ofour/four-ofour.component';
import {ScreenService} from './screen.service';
import {ActionModule} from '../action/action.module';
import {ContentModule} from '../content/content.module';
import {UtilsModule} from './utils.module';
import {MediaModule} from '../media/media.module';
import {MailMessageModule} from '../mailmessage/mailmessage.module';
import {NavigationModule} from '../navigation/navigation.module';
import {PlaceModule} from '../place/place.module';
import {TypedefService} from './typedef.service';
import {SubscriptionModule} from '../subscription/subscription.module';

export function startupServiceFactory(app: Spellbook) {
	const result = () => app.init();
	return result;
}

@NgModule({
	declarations: [
		SpellbookComponent,
		LoginComponent,
		FourOFourComponent,
	],
	imports: [
		MatCardModule,
		MatInputModule,
		MatFormFieldModule,
		MatButtonModule,
		MatToolbarModule,
		MatSidenavModule,
		MatListModule,
		MatIconModule,
		MatMenuModule,
		MatSelectModule,
		ReactiveFormsModule,
		FormsModule,
		CommonModule,
		HttpClientModule,
		FlexLayoutModule,
		SpellbookRoutingModule,
		UsersModule,
		ActionModule,
		ContentModule,
		MediaModule,
		MailMessageModule,
		NavigationModule,
		UtilsModule,
		PlaceModule,
		SubscriptionModule
		// RouterModule.forChild(SpellbookModule.routes),
	],
	exports: [
		SpellbookComponent,
	],
})

export class SpellbookModule {
	// static readonly routes: Routes = [
	// 	{
	// 		path: '',
	// 		redirectTo: 'content',
	// 		pathMatch: 'full'
	// 	},
	// 	{
	// 		path: 'users',
	// 		children: UsersModule.routes
	// 	},
	// 	{
	// 		path: 'content',
	// 		children: ContentModule.routes
	// 	},
	// 	{
	// 		path: 'mailmessage',
	// 		children: MailMessageModule.routes
	// 	},
	// 	{
	// 		path: 'media',
	// 		children: MediaModule.routes
	// 	},
	// 	{
	// 		path: 'place',
	// 		children: PlaceModule.routes
	// 	},
	// 	{
	// 		path: 'page',
	// 		children: NavigationModule.routes
	// 	},
	// 	{
	// 		path: 'action',
	// 		children: ActionModule.routes
	// 	}
	// ];

	static forRoot(config: SpellbookConfig): ModuleWithProviders {
		const providers = PlaceModule.forRoot(config.googleMapKey).providers;
		return {
			ngModule: SpellbookModule,
			providers: [
				providers,
				Spellbook,
				AuthService,
				ScreenService,
				TypedefService,
				{
					provide: APP_INITIALIZER,
					useFactory: startupServiceFactory,
					deps: [Spellbook],
					multi: true
				},
				{
					provide: 'config',
					useValue: config
				},
				{provide: MAT_SNACK_BAR_DEFAULT_OPTIONS, useValue: {duration: 2000}}
			]
		};
	}
}
