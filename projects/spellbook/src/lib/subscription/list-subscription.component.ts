import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import {MatPaginator, MatSnackBar, MatSort, Sort} from '@angular/material';
import {Subscription} from './subscription';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';
import {merge, of as observableOf} from 'rxjs';
import {SubscriptionClient} from './subscription-client';
import {Spellbook} from '../core/spellbook';
import {Client, Filter, ListResponse} from '../core/client';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {MailMessage} from '../mailmessage/mailmessage';

@Component({
	selector: 'splbk-list-mailmessage',
	templateUrl: './list-subscription.component.html',
	styleUrls: ['./list-subscription.component.css']
})
export class ListSubscriptionComponent implements OnInit, AfterViewInit {

	private readonly client: SubscriptionClient;
	public columns: string[];
	public subscribes: Subscription[];
	public currentSize: number;
	public error: string;

	public isLoading: boolean;
	public isLoadingCsv: boolean;

	@ViewChild(MatPaginator) paginator: MatPaginator;
	@ViewChild(MatSort) sort: MatSort;

	private order: string;
	private orderField: string;

	constructor(private spellbook: Spellbook, private snackBar: MatSnackBar) {
		this.columns = [
			'email',
			'country',
			'firstName',
			'lastName',
			'organization',
			'position',
			'created',
			'delete'
		];
		this.client = new SubscriptionClient(spellbook);
		this.subscribes = [];
		this.isLoading = true;
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
					return this.client.getSubscriptionList(this.paginator.pageIndex, this.paginator.pageSize, filters, this.orderField, this.order);
				}),
				map((response: ListResponse<Subscription>) => {
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
			(posts: Subscription[]) => {
				this.subscribes = posts;
			});
	}

	public selectRow(content: Subscription) {
		this.spellbook.router.navigate(['/subscription', content.key]);
	}

	delete(subscription: Subscription) {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + subscription.email + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			snackBarRef.dismiss();
			this.error = undefined;
			this.client.deleteSubscription(subscription).subscribe(
				() => {
					this.snackBar.open('Subscribe deleted!', 'ok', {});
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

	ngOnInit(): void {
	}

	downloadCSV() {
		this.isLoadingCsv = true;
		this.client.downloadCsv().subscribe((result: boolean) => {
			console.log('ok');
			this.isLoadingCsv = false;
		}, (error) => {
			this.isLoadingCsv = false;
		});
	}


}
