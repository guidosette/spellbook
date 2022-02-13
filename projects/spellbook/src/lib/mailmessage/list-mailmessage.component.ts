import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';

import {MatPaginator, MatSnackBar, MatSort, Sort} from '@angular/material';
import {MailMessage} from './mailmessage';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';
import {merge, of as observableOf} from 'rxjs';

import {MessageClient} from './message-client';
import {Spellbook} from '../core/spellbook';
import {Client, Filter, ListResponse} from '../core/client';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {Content} from '../content/content';

@Component({
	selector: 'splbk-list-mailmessage',
	templateUrl: './list-mailmessage.component.html',
	styleUrls: ['./list-mailmessage.component.css']
})
export class ListMailMessageComponent implements OnInit, AfterViewInit {

	private readonly client: MessageClient;
	public columns: string[];
	public mailMessages: MailMessage[];
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
			'sender',
			'recipient',
			'object',
			'created',
			'delete'
		];
		this.client = new MessageClient(spellbook);
		this.mailMessages = [];
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
					return this.client.getMailMessageList(this.paginator.pageIndex, this.paginator.pageSize, filters, this.orderField, this.order);
				}),
				map((response: ListResponse<MailMessage>) => {
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
			(posts: MailMessage[]) => {
				this.mailMessages = posts;
			});
	}

	public selectRow(content: MailMessage) {
		this.spellbook.router.navigate(['/mailmessage', content.id]);
	}

	delete(mailmessage: MailMessage) {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + mailmessage.recipient + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			console.log('actionOk');
			snackBarRef.dismiss();
			this.error = undefined;
			this.client.deleteMailMessage(mailmessage).subscribe(
				() => {
					this.snackBar.open('MailMessage deleted!', 'ok', {});
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

	ngOnInit(): void {
	}

}
