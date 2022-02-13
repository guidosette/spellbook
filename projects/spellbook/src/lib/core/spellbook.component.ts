import {AfterViewInit, Component, Input, ViewChild} from '@angular/core';
import {Spellbook} from './spellbook';
import {User} from './user';
import {Menu} from './menu';
import {ScreenService, ScreenSize} from './screen.service';
import {MatSidenav, MatSnackBar} from '@angular/material';
import {Category} from './category';
import {environment} from '../environments/environment';
import {HttpErrorResponse} from '@angular/common/http';

@Component({
	selector: 'splbk-root',
	templateUrl: './spellbook.component.html',
	styleUrls: ['./spellbook.component.scss']
})
export class SpellbookComponent implements AfterViewInit {

	@ViewChild('sidenav') sidenav: MatSidenav;
	@Input() menu: Menu[];

	public version: string;
	public screenWidth: number;
	public screenWidthLimit = 991;

	constructor(public spellbook: Spellbook, public screenService: ScreenService, public snackBar: MatSnackBar) {
		this.menu = [];

		this.screenService.screenSize.asObservable().subscribe((screenSize: ScreenSize) => {
			this.screenWidth = screenSize.width;
		});

		this.version = environment.version;
		this.spellbook.identity.subscribe(
			(user: User) => {
				this.buildMenu(user);
			}
		);
	}

	ngAfterViewInit(): void {
		setTimeout(() => {
			if (this.sidenav) {
				this.sidenav.open();
			}
		}, 200);
	}

	protected buildMenu(user: User) {
		if (!user) {
			this.menu = [];
			return;
		}

		if (user.hasPermission(User.PERMISSION_READ_USER) || user.hasPermission(User.PERMISSION_WRITE_USER)) {
			const users = new Menu('users', '/users', 'Users', 'supervised_user_circle');
			this.menu.push(users);
		}
		if (user.hasPermission(User.PERMISSION_READ_CONTENT) || user.hasPermission(User.PERMISSION_WRITE_CONTENT)) {
			const types: Array<Menu> = new Array<Menu>();
			this.spellbook.getSupportedCategories().forEach((category: Category) => {
				const c = new Menu(category.name, '/content', category.label, null, {
					type: category.type,
					category: category.name
				});
				// get the supermenu if it exists, create it otherwise
				let parent = types.find(m => m.id === category.type);
				if (!parent) {
					let label = category.type;
					let icon = '';
					const def = this.spellbook.definitions.getTypeDefinition(category.type);
					if (def) {
						label = def.menuLabel();
						icon = def.menuIcon();
					}
					parent = new Menu(category.type, null, label, icon);
					types.push(parent);
				}
				parent.addChildren(c);
			});

			for (const t of types) {
				if (t.children.length > 0) {
					this.menu.push(t);
				}
			}
		}
		if (user.hasPermission(User.PERMISSION_READ_MAILMESSAGE) || user.hasPermission(User.PERMISSION_WRITE_MAILMESSAGE)) {
			const mailmessage = new Menu('mailmessage', '/mailmessage', 'Mail Messages', 'mail_outline');
			this.menu.push(mailmessage);
		}
		if (user.hasPermission(User.PERMISSION_READ_SUBSCRIPTION) || user.hasPermission(User.PERMISSION_WRITE_SUBSCRIPTION)) {
			const subscription = new Menu('subscription', '/subscription', 'Subscription', 'subscriptions', null);
			this.menu.push(subscription);
		}
		if (user.hasPermission(User.PERMISSION_READ_MEDIA) || user.hasPermission(User.PERMISSION_WRITE_MEDIA)) {
			const media = new Menu('media', '/media', 'Media', 'perm_media');
			this.menu.push(media);
		}
		if (user.hasPermission(User.PERMISSION_READ_PLACE) || user.hasPermission(User.PERMISSION_WRITE_PLACE)) {
			const place = new Menu('place', '/place', 'Places', 'place');
			this.menu.push(place);
		}
		if (user.hasPermission(User.PERMISSION_READ_PAGE) || user.hasPermission(User.PERMISSION_WRITE_PAGE)) {
			const pageMenu = new Menu('page', '/page', 'Pages', 'bookmark_border');
			this.menu.push(pageMenu);
		}

		if (user.hasPermission(User.PERMISSION_READ_ACTION) || user.hasPermission(User.PERMISSION_WRITE_ACTION)) {
			const actions = new Menu('action', '/action', 'Actions', 'build', null);
			this.menu.push(actions);
		}
	}

	public logout() {
		this.spellbook.token = null;
		// navigate to the login
		this.spellbook.router.navigate(['/login']);
	}
}
