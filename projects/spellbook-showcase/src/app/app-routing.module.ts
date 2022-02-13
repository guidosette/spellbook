import {NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {UsersModule} from '../../../spellbook/src/lib/users/users.module';
import {TestComponent} from './test.component';
import {NavigationModule} from '../../../spellbook/src/lib/navigation/navigation.module';
import {PlaceModule} from '../../../spellbook/src/lib/place/place.module';
import {MailMessageModule} from '../../../spellbook/src/lib/mailmessage/mailmessage.module';
import {ContentModule} from '../../../spellbook/src/lib/content/content.module';
import {MediaModule} from '../../../spellbook/src/lib/media/media.module';
import {TestModule} from './test/test.module';
import {ActionModule} from '../../../spellbook/src/lib/action/action.module';
import {SubscriptionModule} from '../../../spellbook/src/lib/subscription/subscription.module';

@NgModule({
	imports: [
		RouterModule.forRoot([
			{
				path: '',
				component: TestComponent, // SpellbookComponent
				children: [
					// {path: 'core', loadChildren: '../../../spellbook/src/lib/core/spellbook.module#SpellbookModule'},
					{
						path: 'users',
						children: UsersModule.routes
					},
					{
						path: 'content',
						children: ContentModule.routes
					},
					{
						path: 'mailmessage',
						children: MailMessageModule.routes
					},
					{
						path: 'media',
						children: MediaModule.routes
					},
					{
						path: 'place',
						children: PlaceModule.routes
					},
					{
						path: 'page',
						children: NavigationModule.routes
					},
					{
						path: 'action',
						children: ActionModule.routes
					},
					{
						path: 'subscription',
						children: SubscriptionModule.routes
					},
					{
						path: 'test',
						children: TestModule.routes
					},
				]
			}
		])
	],
	exports: [RouterModule]
})
export class AppRoutingModule {

}
