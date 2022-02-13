import {Spellbook} from '../../../spellbook/src/lib/core/spellbook';
import {ScreenService} from '../../../spellbook/src/lib/core/screen.service';
import {Component} from '@angular/core';
import {User} from '../../../spellbook/src/lib/core/user';
import {Menu} from '../../../spellbook/src/lib/core/menu';
import {SpellbookComponent} from '../../../spellbook/src/lib/core/spellbook.component';
import {MatSnackBar} from '@angular/material';

@Component({
	selector: 'test-root',
	templateUrl: './test.component.html',
	styleUrls: ['./test.component.scss']
})
export class TestComponent extends SpellbookComponent {

	constructor(public spellbook: Spellbook, public screenService: ScreenService, public snackBar: MatSnackBar) {
		super(spellbook, screenService, snackBar);

	}

	protected buildMenu(user: User) {
		super.buildMenu(user);
		const test = new Menu('test', '/test', 'Test', 'bug_report', null);
		this.menu.push(test);
	}
}
