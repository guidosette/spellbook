import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import {Action} from '../action/action';
import {MatDialog, MatPaginator, MatSnackBar, MatSort, Sort} from '@angular/material';
import {Spellbook} from '../core/spellbook';
import {Client, Filter, ListResponse} from '../core/client';
import {merge, Observable, of as observableOf} from 'rxjs';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';
import {ActionClient} from './action-client';
import {AttachmentType} from '../core/supported-attachment';
import {Attachment} from '../media/multimedia/attachment';
import {UploadMediaFileComponent} from '../media/upload-media-file.component';
import {FileAttachment} from '../media/multimedia/file';
import {Task} from './task';
import {HttpErrorResponse} from '@angular/common/http';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';

@Component({
	selector: 'splbk-list-action',
	templateUrl: './list-action.component.html',
	styleUrls: ['./list-action.component.scss']
})
export class ListActionComponent implements OnInit, AfterViewInit {
	private readonly client: ActionClient;

	public columns: string[];
	public actions: Action[];
	public currentSize: number;
	public error: string;

	public isLoading;
	public isLoadingTask;

	private order: string;
	private orderField: string;

	@ViewChild(MatPaginator) paginator: MatPaginator;
	@ViewChild(MatSort) sort: MatSort;

	constructor(private spellbook: Spellbook, private snackBar: MatSnackBar, public dialog: MatDialog) {
		this.columns = [
			'name',
			'endpoint',
			'method',
			'type',
			'created',
			'status',
			'count',
			'message',
			'action',
		];
		this.client = new ActionClient(this.spellbook);

		this.actions = [];
		this.isLoading = true;

		this.orderField = 'Name';
		this.order = Client.orderAscKey;

	}

	ngOnInit(): void {
		this.load(0);
	}

	ngAfterViewInit() {
		this.sort.sortChange.subscribe(() => {
			this.paginator.pageIndex = 0;
		});

		this.sortData({
			active: 'Name',
			direction: 'asc'
		});

	}

	getTaskList() {
		this.isLoadingTask = true;
		this.client.getTasks().subscribe((response: ListResponse<Task>) => {
			this.isLoadingTask = false;
			this.setTasksToActions(response.items);
		}, (err: HttpErrorResponse) => {
			this.isLoadingTask = false;
			this.showError(err.error.Error);
		});
	}

	setTasksToActions(tasks: Array<Task>) {
		this.actions.forEach((action: Action) => {
			const index = tasks.findIndex((task: Task) => {
				return task.url === action.endpoint;
			});
			if (index >= 0) {
				action.task = new Task(tasks[index]);
			} else {
				action.task = null;
			}
		});
	}

	public showError(message: string) {
		console.error('error:', message);
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = message;
		snackbarData.isError = true;
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			snackBarRef.dismiss();
		};
	}

	sortData(sort: Sort) {
		if (!sort.active || sort.direction === '') {
			return;
		}
		this.orderField = sort.active;
		this.order = sort.direction === 'asc' ? Client.orderAscKey : Client.orderDescKey;
	}

	public load(index: number) {
		this.getTaskList();

		this.paginator.pageIndex = index;
		merge(this.sort.sortChange, this.paginator.page)
			.pipe(
				startWith({}),
				switchMap(() => {
					this.isLoading = true;
					const filters: Filter[] = [];
					return this.client.getActionList(this.paginator.pageIndex, this.paginator.pageSize, filters, this.orderField, this.order);
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
			(actions: Action[]) => {
				this.actions = actions;
			});
	}

	public selectRow(action: Action) {
		// this.spellbook.router.navigate(['/action', action.endpoint]);
	}

	runAction(action: Action) {
		if (action.isRun()) {
			action.status = Action.ACTION_STATUS_STOP;
		} else {
			if (action.type === Action.ACTION_TYPE_UPLOAD) {
				this.upload(action).subscribe(() => {
					this.createRunTask(action);
				});
			} else {
				this.createRunTask(action);
			}
		}
	}

	public createRunTask(action: Action): void {
		if (!action.isTaskCreated()) {
			this.createTask(action).subscribe((task: Task) => {
				action.task = new Task(task);
				this.runTask(action);
			});
		} else {
			this.runTask(action);
		}
	}

	public createTask(action: Action): Observable<Task> {
		return new Observable<Task>((ob) => {
			action.status = Action.ACTION_STATUS_RUN;
			this.client.createTask(action.getTask()).subscribe(
				(task: Task) => {
					console.log('task', task);
					ob.next(task);
				},
				(err: HttpErrorResponse) => {
					action.status = Action.ACTION_STATUS_STOP;
					this.showError(err.error.Error);
					ob.error();
				}
			);
		});

	}

	public runTask(action: Action): void {
		this.client.runTask(action.getTask()).subscribe(
			(task: Task) => {
				console.log('task', task);
				action.status = Action.ACTION_STATUS_RUN;
			},
			(err: HttpErrorResponse) => {
				action.status = Action.ACTION_STATUS_STOP;
				this.showError(err.error.Error);
			}
		);
	}

	public upload(action: Action): Observable<FileAttachment> {
		return new Observable<FileAttachment>((ob) => {
			const attachment = new Attachment();
			attachment.group = 'Excel';
			attachment.type = AttachmentType.ATTACHMENT;
			attachment.parentKey = action.endpoint;
			const dialogRef = this.dialog.open(UploadMediaFileComponent, {
				width: '300px',
				data: attachment
			});
			dialogRef.componentInstance.multipleMode = false;
			dialogRef.componentInstance.folder = 'upload';
			dialogRef.componentInstance.namespace = 'actions';
			dialogRef.afterClosed().subscribe((result: FileAttachment) => {
				if (!result) {
					return;
				}
				console.log('results', result);
				// this.snackBar.open('File uploaded!', 'ok', {});
				ob.next(result);
			});
		});
	}
}
