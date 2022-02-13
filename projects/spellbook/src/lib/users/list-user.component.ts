import {AfterViewInit, Component, ViewChild} from '@angular/core';

import {MatPaginator, MatSort, Sort} from '@angular/material';
import {merge, of as observableOf} from 'rxjs';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';
import {User} from '../core/user';
import {Spellbook} from '../core/spellbook';
import {Client, Filter} from '../core/client';


@Component({
	selector: 'splbk-list-user',
	templateUrl: './list-user.component.html',
	styleUrls: ['./list-user.component.css']
})
export class ListUserComponent implements AfterViewInit {

	public columns: string[];
	public users: User[];
	public currentSize: number;

	public isLoading;

	public create: boolean;

	@ViewChild(MatPaginator) paginator: MatPaginator;
	@ViewChild(MatSort) sort: MatSort;

	private order: string;
	private orderField: string;

	constructor(private spellbook: Spellbook) {
		this.columns = [
			'username',
			'name',
			'surname',
			'email'
		];
		this.users = [];
		this.isLoading = true;
	}

	ngAfterViewInit() {
		this.sort.sortChange.subscribe(() => {
			this.paginator.pageIndex = 0;
		});

		this.sortData({
			active: 'Email',
			direction: 'desc'
		});
		this.load(0);
	}

	sortData(sort: Sort) {
		if (!sort.active || sort.direction === '') {
			return;
		}
		this.orderField = sort.active;
		this.order = sort.direction === 'asc' ? Client.orderAscKey : Client.orderDescKey;
	}

	public load(index: number) {
		this.paginator.pageIndex = index;
		merge(this.sort.sortChange, this.paginator.page)
			.pipe(
				startWith({}),
				switchMap(() => {
					this.isLoading = true;
					const filters: Filter[] = [];
					return this.spellbook.api.getUsers(this.paginator.pageIndex, this.paginator.pageSize, filters, this.orderField, this.order);
				}),
				map((response: any) => {
					this.isLoading = false;
					const items = response.items;
					const size = (this.paginator.pageIndex + 1) * Math.min(this.paginator.pageSize, items.length);
					this.currentSize = response.more ? size + 1 : size;
					return items;
				}),
				catchError(() => {
					this.isLoading = false;
					return observableOf([]);
				})
			).subscribe(
			(users: User[]) => {
				this.users = users;
			});
	}

	public selectRow(user: User) {
		this.spellbook.router.navigate([`/users/${user.username}`]);
	}

	public toggleMode() {
		this.create = !this.create;
	}

}
