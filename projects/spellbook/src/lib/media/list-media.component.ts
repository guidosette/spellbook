import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import {Attachment} from './multimedia/attachment';
import {Spellbook} from '../core/spellbook';
import {ActivatedRoute} from '@angular/router';
import {MatDialog, MatDialogRef, MatPaginator} from '@angular/material';
import {ScreenService, ScreenSize} from '../core/screen.service';
import {CreateUpdateMediaComponent} from './create-update-media.component';
import {BrowseMediaFilesComponent} from './browse-media-files.component';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';
import {of as observableOf} from 'rxjs';
import {UploadMediaFileComponent} from './upload-media-file.component';
import {CollectMediaUrlComponent} from './collect-media-url.component';
import {Client, Filter} from '../core/client';
import {MediaClient} from './media-client';
import {FileAttachment} from './multimedia/file';
import {UrlUtils} from '../core/url-utils';

export class DropDownItem {
	value: string;
	viewValue: string;

	constructor(value: string, viewValue: string) {
		this.value = value;
		this.viewValue = viewValue;
	}
}

@Component({
	selector: 'splbk-list-media',
	templateUrl: './list-media.component.html',
	styleUrls: ['./list-media.component.scss']
})
export class ListMediaComponent implements OnInit, AfterViewInit {

	private static VALUE_ANY: DropDownItem = {value: null, viewValue: 'Any'};
	static MINE_LABEL = 'Mine';

	// -------------------------------filterType start
	public set filterType(value: DropDownItem) {
		this._filterType = value;
		// this.refreshData(0);
	}

	public get filterType() {
		return this._filterType;
	}

	private _filterType: DropDownItem = ListMediaComponent.VALUE_ANY;
	public attachmentTypes: DropDownItem[];
	// -------------------------------filterType start

	// -------------------------------filterGroup start
	public set filterGroup(value: DropDownItem) {
		this._filterGroup = value;
		// this.refreshData(0);
	}

	public get filterGroup() {
		return this._filterGroup;
	}

	private _filterGroup: DropDownItem = ListMediaComponent.VALUE_ANY;
	public attachmentGroups: DropDownItem[];
	// -------------------------------filterGroup end

	// -------------------------------filterParent start
	public set filterParent(value: DropDownItem) {
		this._filterParent = value;
		// this.refreshData(0);
	}

	public get filterParent() {
		return this._filterParent;
	}

	private _filterParent: DropDownItem = ListMediaComponent.VALUE_ANY;
	public attachmentParents: DropDownItem[];
	// -------------------------------filterParent end

	// DATA
	public type: string = Attachment.TYPE_DEFAULT;
	public group: string = Attachment.GROUP_DEFAULT;
	public parentKey: string = Attachment.PARENT_DEFAULT;

	constructor(protected spellbook: Spellbook, private route: ActivatedRoute, private screenService: ScreenService,
	            public dialog: MatDialog) {
		this.client = new MediaClient(spellbook);
		this.isLoading = true;
		this.screenService.screenSize.asObservable().subscribe((screenSize: ScreenSize) => {
			this.screenWidth = screenSize.width;
			this.setColGrid();
		});
		this.multipleMode = true;
	}

	protected readonly client: MediaClient;

	public filterableByType = true;
	public filterableByGroup = false;
	public filterableByParent = false;
	public media: Attachment[];
	public screenWidth: number;
	public colspanGrid: number;
	public secondaryFabVisible = false;
	public currentSize: number;
	public isLoading;

	// type of the parent to which the attachment relates
	public parentType: string;

	public multipleMode: boolean;

	public utils = UrlUtils;

	@ViewChild(MatPaginator) paginator: MatPaginator;

	private static arrayToKeyValue(values: string[]) {
		const kv: DropDownItem[] = [
			ListMediaComponent.VALUE_ANY,     // Default unselected value
		];
		for (const v of values) {
			const item: DropDownItem = {
				value: v,
				viewValue: v
			};
			kv.push(item);
		}
		return kv;
	}

	// Creates and populates a new attachment with context-specific values
	private buildAttachment(withParent?: boolean): Attachment {
		const attachment = new Attachment();
		attachment.group = this.group;
		attachment.type = this.filterType.value;
		if (withParent) {
			attachment.parentKey = this.parentKey;
		} else {
			attachment.parentKey = Attachment.PARENT_DEFAULT;
			attachment.group = Attachment.GROUP_DEFAULT;
		}
		attachment.parentType = this.parentType;
		return attachment;
	}

	public setColGrid() {
		if (this.screenWidth <= 991) {
			this.colspanGrid = 2;
		} else if (this.screenWidth <= 470) {
			this.colspanGrid = 3;
		}
		this.colspanGrid = 1;
	}

	ngOnInit() {
		this.getTypes();
		this.refreshSelects();
	}

	refreshSelects() {
		this.getGroups();
		this.getParents();
	}

	ngAfterViewInit() {
		this.refreshData(0);
	}


	private getGroups() {
		// Download Group data
		this.client.getAttachmentPropertyList('Group').subscribe(
			(groups: string[]) => {
				this.attachmentGroups = ListMediaComponent.arrayToKeyValue(groups);
				this.selectDefaultGroup();
			},
			(error) => {
				console.log(`Unable to load group data due to error ${error.toString()}`);
			}
		);
	}

	private selectDefaultGroup() {
		if (this.filterGroup) {
			this.filterGroup = this.attachmentGroups.find((dd: DropDownItem) => {
				return dd.value === this.filterGroup.value;
			});
		}
	}

	private getParents() {
		// Download Parent data
		if (this.parentKey) {
			const kv: DropDownItem[] = [
				ListMediaComponent.VALUE_ANY,     // Default unselected value
				{value: this.parentKey, viewValue: ListMediaComponent.MINE_LABEL},     // Mine parent
			];
			this.attachmentParents = kv;
			this.selectDefaultParent();
			return;
		}
		this.client.getAttachmentPropertyList('ParentKey').subscribe(
			(parents: string[]) => {
				this.attachmentParents = ListMediaComponent.arrayToKeyValue(parents);
				this.selectDefaultParent();
			},
			(error) => {
				console.log(`Unable to load parent data due to error ${error.toString()}`);
			}
		);
	}

	private selectDefaultParent() {
		if (this.filterParent) {
			this.filterParent = this.attachmentParents.find((dd: DropDownItem) => {
				return dd.value === this.filterParent.value;
			});
		}
	}

	private getTypes() {
		this.attachmentTypes = [
			ListMediaComponent.VALUE_ANY     // Default unselected value
		];
		for (const v of this.spellbook.supportedAttachments) {
			const item: DropDownItem = {
				value: v.value,
				viewValue: v.name
			};
			this.attachmentTypes.push(item);
		}
		this.selectDefaultType();
	}

	private selectDefaultType() {
		if (this.filterType) {
			this.filterType = this.attachmentTypes.find((dd: DropDownItem) => {
				return dd.value === this.filterType.value;
			});
		}
	}

	public refreshData(index: number) {
		if (!this.paginator) {
			// Not yet initalized
			return;
		}
		// Download Media data
		this.paginator.pageIndex = index;
		this.paginator.page
			.pipe(
				startWith({}),
				switchMap(() => {
					this.isLoading = true;

					const filters: Filter[] = [];
					if (this.filterType.value) {
						filters.push(new Filter('Type', this.filterType.value));
					}
					if (this.filterGroup.value) {
						filters.push(new Filter('Group', this.filterGroup.value));
					}
					if (this.filterParent.value) {
						filters.push(new Filter('ParentKey', this.filterParent.value));
					}
					return this.client.getAttachmentList(this.paginator.pageIndex, this.paginator.pageSize, filters, 'Created', Client.orderDescKey);
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
			(a: Attachment[]) => {
				this.media = a;
			});
	}

	onItemClick(item: Attachment, forceCreate?: boolean) {
		this.showDetail(item);
	}

	showDetail(item: Attachment) {
		// Gallery: open delected item's detail
		const dialogRef = this.dialog.open(CreateUpdateMediaComponent, {
			width: '800px',
			data: item
		});
		dialogRef.componentInstance.showImageFullScreen = true;

		dialogRef.afterClosed().subscribe((a: Attachment) => {
			setTimeout(() => {
				if (a) {
					this.refreshData(this.paginator.pageIndex);
					this.refreshSelects();
				}
			}, 0);
		});
	}

	showDetails(items: Attachment[]) {
		items.forEach((item: Attachment) => {
			this.showDetail(item);
		});
	}

	addFromLocalFile() {
		this.secondaryFabVisible = false;
		const dialogRef = this.dialog.open(UploadMediaFileComponent, {
			width: '300px',
			data: this.buildAttachment()
		});
		dialogRef.componentInstance.multipleMode = this.multipleMode;
		this.openDialog(dialogRef);
	}

	addFromUrl() {
		this.secondaryFabVisible = false;
		const dialogRef = this.dialog.open(CollectMediaUrlComponent, {
			width: '600px',
			data: this.buildAttachment()
		});
		this.openDialog(dialogRef);
	}

	addFromServer() {
		this.secondaryFabVisible = false;
		const dialogRef = this.dialog.open(BrowseMediaFilesComponent, {
			width: '800px',
			height: '90vh',
			data: this.buildAttachment()
		});
		this.openDialog(dialogRef);
	}

	private openDialog(dialogRef: MatDialogRef<any>) {
		dialogRef.afterClosed().subscribe((results: FileAttachment[]) => {
			// results: FileAttachment[] o FileAttachment
			if (!results) {
				return;
			}
			if (results instanceof Array) {
				console.log('results multiple', results);
				// multiple
				const newAttachments: Attachment[] = [];
				results.forEach((result: FileAttachment) => {
					const newAttachment = this.buildAttachment(true);
					newAttachment.name = result.name;
					newAttachment.resourceUrl = result.resourceUrl;
					newAttachment.resourceThumbUrl = result.resourceThumbUrl;
					newAttachment.type = result.attachmentType;
					newAttachments.push(newAttachment);
				});
				this.showDetails(newAttachments);
			} else {
				console.log('results single', results);
				// single
				const result: FileAttachment = results;
				const newAttachment = this.buildAttachment(true);
				newAttachment.name = result.name;
				newAttachment.resourceUrl = result.resourceUrl;
				newAttachment.resourceThumbUrl = result.resourceThumbUrl;
				newAttachment.type = result.attachmentType;
				console.log('newAttachment', newAttachment);
				this.showDetail(newAttachment);
			}
		});
	}

	public getIconForAttachmentName(name: string) {
		return Attachment.getIconForAttachmentName(name);
	}

	public resetFilters() {
		if (this.filterableByType) {
			this.filterType = ListMediaComponent.VALUE_ANY;
		}
		if (this.filterableByGroup) {
			this.filterGroup = ListMediaComponent.VALUE_ANY;
		}
		if (this.filterableByParent) {
			this.filterParent = ListMediaComponent.VALUE_ANY;
		}
		this.refreshData(0);
	}

	getStyle(attachment: Attachment) {
		if (this.parentKey === attachment.parentKey && this.group === attachment.group && this.parentKey !== Attachment.PARENT_DEFAULT) {
			// same group and parent
			return {
				border: 'solid 4px #9c27b0',
				height: 'calc(100% - 8px)',
				width: '100%',
			};
		}
		return {
			border: 'solid 0px #9c27b0',
			height: '100%',
			width: '100%',
		};
	}


}
