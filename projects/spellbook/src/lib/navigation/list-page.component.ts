import {AfterViewInit, Component, ViewChild} from '@angular/core';
import {Spellbook} from '../core/spellbook';
import {MatPaginator, MatSnackBar, MatSort, Sort} from '@angular/material';
import {ActivatedRoute} from '@angular/router';
import {Page} from './page';
import {Client, Filter, ListResponse} from '../core/client';
import {merge, of as observableOf} from 'rxjs';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';
import {CreateUpdatePageComponent} from './create-update-page.component';
import {PlatformLocation} from '@angular/common';
import {NavigationClient} from './navigation-client';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';

@Component({
	selector: 'splbk-list-page',
	templateUrl: './list-page.component.html',
	styleUrls: ['./list-page.component.css']
})
export class ListPageComponent implements AfterViewInit {
	private client: NavigationClient;

	public columns: string[];
	public pages: Page[];
	public currentSize: number;
	public isLoading;
	public create: boolean;
	private order: string;
	private orderField: string;
	public error: string;
	public baseUrl: string;

	@ViewChild(MatPaginator) paginator: MatPaginator;
	@ViewChild(MatSort) sort: MatSort;
	@ViewChild(CreateUpdatePageComponent) createUpdatePageComponent: CreateUpdatePageComponent;

	public allLanguages: string[];

	public set languageSelected(value: string) {
		this._languageSelected = value;
		this.load(0);
	}

	public get languageSelected(): string {
		return this._languageSelected;
	}

	private _languageSelected: string;

	private allLanguageLabel = 'All';

	constructor(private spellbook: Spellbook, private snackBar: MatSnackBar, private route: ActivatedRoute, private platformLocation: PlatformLocation) {
		this.baseUrl = (platformLocation as any).location.origin;
		this.columns = [
			'URL',
			'label',
			'code',
			'locale',
			'delete'
		];
		this.pages = [];
		this.isLoading = true;
		this.client = new NavigationClient(spellbook);

		this.orderField = 'Title';
		this.order = Client.orderAscKey;

		this.setAllLanguages();

	}

	ngAfterViewInit(): void {
		this.sort.sortChange.subscribe(() => {
			this.paginator.pageIndex = 0;
		});

		this.sortData({
			active: 'Title',
			direction: 'asc'
		});
	}

	private setAllLanguages() {
		this.spellbook.api.getLanguages().subscribe((allLanguages: string[]) => {
			this.allLanguages = allLanguages;
			this.allLanguages.unshift(this.allLanguageLabel);
			this.languageSelected = this.allLanguages.length > 0 ? this.allLanguages[1] : this.allLanguageLabel;
		});
	}

	public load(index: number) {
		this.paginator.pageIndex = index;
		merge(this.sort.sortChange, this.paginator.page)
			.pipe(
				startWith({}),
				switchMap(() => {
					this.isLoading = true;
					const filters: Filter[] = [];
					if (this.languageSelected && this.languageSelected !== this.allLanguageLabel) {
						filters.push(new Filter('Locale', this.languageSelected));
					}
					return this.client.getPageList(this.paginator.pageIndex, this.paginator.pageSize, filters, this.orderField, this.order);
				}),
				map((response: ListResponse<Page>) => {
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
			(pages: Page[]) => {
				this.pages = pages;
			});
	}

	public toggleMode() {
		this.create = !this.create;
	}

	delete(p: Page) {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete page with title ' + p.title + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			this.error = undefined;
			snackBarRef.dismiss();
			this.client.deletePage(p).subscribe(
				() => {
					this.snackBar.open('Page removed', 'ok', {});
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

	public selectRow(p: Page) {
		this.spellbook.router.navigate(['/page', p.id]);
	}

	sortData(sort: Sort) {
		if (!sort.active || sort.direction === '') {
			return;
		}
		this.orderField = sort.active;
		this.order = sort.direction === 'asc' ? Client.orderAscKey : Client.orderDescKey;
	}

}
