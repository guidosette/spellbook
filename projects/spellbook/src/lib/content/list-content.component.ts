import {AfterViewInit, Component, ViewChild} from '@angular/core';

import {MatPaginator, MatSnackBar, MatSort, Sort} from '@angular/material';
import {merge, of as observableOf} from 'rxjs';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';
import {Content} from './content';

import {CreateUpdateContentComponent} from './create-update-content.component';

import {ActivatedRoute} from '@angular/router';
import {Category} from '../core/category';
import {ContentClient} from './content-client';
import {Spellbook} from '../core/spellbook';
import {Client, Filter, ListResponse} from '../core/client';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {ContentDefinition} from './content-definition';
import {UrlUtils} from '../core/url-utils';

@Component({
	selector: 'splbk-list-post',
	templateUrl: './list-content.component.html',
	styleUrls: ['./list-content.component.css']
})
export class ListContentComponent implements AfterViewInit {
	private readonly client: ContentClient;

	public columns: string[];
	public posts: Content[];
	public currentSize: number;
	public error: string;

	public isLoading;

	public create: boolean;
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
	private order: string;
	private orderField: string;

	public currentType: string;
	private currentCategory: string;

	public definition: ContentDefinition;
	public localeHidden: boolean;
	public utils = UrlUtils;

	@ViewChild(MatPaginator) paginator: MatPaginator;
	@ViewChild(MatSort) sort: MatSort;
	@ViewChild(CreateUpdateContentComponent) createUpdateContentComponent: CreateUpdateContentComponent;

	constructor(private spellbook: Spellbook, private snackBar: MatSnackBar, private route: ActivatedRoute) {
		this.columns = [
			'slug',
			'title',
			'locale',
			'code',
			'created',
			'published',
			'order',
			'isPublished',
			// 'cover',
			'delete'
		];

		this.client = new ContentClient(this.spellbook);

		this.posts = [];
		this.isLoading = true;
		this.setAllLanguages();

		this.route.queryParams.subscribe(params => {
			// console.log('queryParams', params);
			this.currentType = params.type;
			this.currentCategory = params.category;

			this.definition = this.spellbook.definitions.getTypeDefinition<Content>(this.currentType) as ContentDefinition;
			// if no definition is found, use the default one
			if (!this.definition) {
				this.definition = new ContentDefinition(this.currentType, '', '');
			}
			this.localeHidden = this.definition.field('locale').hidden;
			this.columns = this.definition.columnFields();
			this.orderField = this.definition.listOrderField();
			this.order = this.definition.listOrder();
			if (this.paginator) {
				this.load(0);
			}
			this.create = false;
		});
	}

	public getCategoryName(): string {
		const cat = this.spellbook.getSupportedCategories().find((s) => {
			return s.name === this.currentCategory;
		});
		return cat ? cat.label : '';
	}

	private setAllLanguages() {
		this.spellbook.api.getLanguages().subscribe((allLanguages: string[]) => {
			this.allLanguages = allLanguages;
			this.allLanguages.unshift(this.allLanguageLabel);
			this.languageSelected = this.allLanguages.length > 0 ? this.allLanguages[1] : this.allLanguageLabel;
		});
	}

	ngAfterViewInit() {
		this.sort.sortChange.subscribe(() => {
			this.paginator.pageIndex = 0;
		});
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
					if (this.languageSelected && this.languageSelected !== this.allLanguageLabel) {
						filters.push(new Filter('Locale', this.languageSelected));
					}
					if (this.currentCategory) {
						filters.push(new Filter('Category', this.currentCategory));
					}
					if (this.currentType) {
						filters.push(new Filter('Type', this.currentType));
					}
					return this.client.getContentList(this.paginator.pageIndex, this.paginator.pageSize, filters, this.orderField, this.order);
				}),
				map((response: ListResponse<Content>) => {
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
			(posts: Content[]) => {
				this.posts = posts;
			});
	}

	public selectRow(content: Content) {
		this.spellbook.router.navigate(['/content', content.id]);
	}

	public toggleMode() {
		this.create = !this.create;
		if (this.create) {
			setTimeout(() => {
				if (this.createUpdateContentComponent) {
					this.createUpdateContentComponent.currentCategory = this.currentCategory;
					this.createUpdateContentComponent.setCurrentType(this.currentType);
				}
			}, 0);
		}
	}

	delete(content: Content) {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + content.title + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			console.log('actionOk');
			this.error = undefined;
			snackBarRef.dismiss();
			this.client.deleteContent(content).subscribe(
				() => {
					this.snackBar.open('Content deleted!', 'ok', {});
					// refresh
					this.load(0);
				},
				(err) => {
					this.error = err.statusText;
				}
			);

		};
		snackbarData.actionNo = () => {
			console.log('actionNo');
			snackBarRef.dismiss();
		};

	}

	clickTranslate(content: Content) {
		// check all translation
		this.client.getContentSupportedLanguages(content.idTranslate).subscribe((languagesInserted: string[]) => {
			this.spellbook.api.getLanguages().subscribe((allLanguages: string[]) => {
				if (allLanguages.length === languagesInserted.length) {
					this.snackBar.open('All translations are complete!', 'ok', {});
				} else {
					this.create = true;
					setTimeout(() => {
						this.createUpdateContentComponent.setTranslateMode(content, allLanguages, languagesInserted);
					});
				}
			});
		});
	}
}
