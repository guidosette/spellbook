import {Component, Input, OnInit} from '@angular/core';
import {Menu} from './menu';

@Component({
	selector: 'splbk-menu-item',
	templateUrl: './menu-item.component.html',
	styleUrls: ['./menu-item.component.scss']
})
export class MenuItemComponent implements OnInit {

	@Input() item: Menu;
	@Input() isChild: boolean;
	@Input() level: number;

	constructor() {
	}

	ngOnInit() {
	}

	setMyStyles() {
		const styles = {
			'margin-left': this.level * 20 + 'px'
		};
		return styles;
	}
}
