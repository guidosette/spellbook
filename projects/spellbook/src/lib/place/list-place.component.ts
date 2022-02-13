import {AfterViewInit, Component, ViewChild} from '@angular/core';
import {Place} from './place';
import {MatPaginator, MatSnackBar, MatSort, Sort} from '@angular/material';
import {Spellbook} from '../core/spellbook';
import {merge, of as observableOf} from 'rxjs';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';

import {ActivatedRoute} from '@angular/router';
import {CreateUpdatePlaceComponent} from './create-update-place.component';
import {PlaceClient} from './place-client';
import {Client, Filter} from '../core/client';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';

@Component({
	selector: 'splbk-list-place',
	templateUrl: './list-place.component.html',
	styleUrls: ['./list-place.component.scss']
})
export class ListPlaceComponent implements AfterViewInit {

	private readonly client: PlaceClient;
	public columns: string[];
	public places: Place[];
	public currentSize: number;
	public error: string;

	public isLoading;
	public searchFilter: Filter;

	public create: boolean;

	private order: string;
	private orderField: string;

	@ViewChild(MatPaginator) paginator: MatPaginator;
	@ViewChild(MatSort) sort: MatSort;
	@ViewChild(CreateUpdatePlaceComponent) createUpdatePlaceComponent: CreateUpdatePlaceComponent;

	constructor(private spellbook: Spellbook, private snackBar: MatSnackBar, private route: ActivatedRoute) {
		this.columns = [
			'name',
			'city',
			'created',
			'delete'
		];
		this.client = new PlaceClient(spellbook);
		this.places = [];
		this.isLoading = true;

		this.orderField = 'Address';
		this.order = Client.orderAscKey;
	}

	ngAfterViewInit() {
		this.sort.sortChange.subscribe(() => {
			this.paginator.pageIndex = 0;
		});

		this.sortData({
			active: 'Created',
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
					if (this.searchFilter) {
						filters.push(this.searchFilter);
					}
					return this.client.getPlaceList(this.paginator.pageIndex, this.paginator.pageSize, filters, this.orderField, this.order);
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
			(places: Place[]) => {
				this.places = places;
			});
	}

	public selectRow(content: Place) {
		this.spellbook.router.navigate(['/place', content.id]);
	}

	public toggleMode() {
		this.create = !this.create;
		if (this.create) {
			setTimeout(() => {
				if (this.createUpdatePlaceComponent) {
					// this.createUpdatePlaceComponent.currentCategory = this.currentCategory;
				}
			}, 0);
		}
	}

	delete(place: Place) {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + place.address + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			this.error = undefined;
			snackBarRef.dismiss();
			this.client.deletePlace(place).subscribe(
				() => {
					this.snackBar.open('Place deleted!', 'ok', {});
					// refresh
					this.load(0);
				},
				(err) => {
					this.error = err.statusText;
				}
			);

		};
		snackbarData.actionNo = () => {
			snackBarRef.dismiss();
		};

	}

}
