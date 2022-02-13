import {Component, Inject} from '@angular/core';
import {ListMediaComponent} from './list-media.component';
import {Spellbook} from '../core/spellbook';
import {ActivatedRoute} from '@angular/router';
import {ScreenService} from '../core/screen.service';
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef, MatSnackBar} from '@angular/material';
import {Attachment} from './multimedia/attachment';
import {SupportedAttachment} from '../core/supported-attachment';
import {forkJoin} from 'rxjs';
import {HttpErrorResponse} from '@angular/common/http';
import {CreateUpdateMediaComponent} from './create-update-media.component';

export interface ListMediaDialogComponentConfig {
	filterableByType: boolean;
	filterableByGroup: boolean;
	filterableByParent: boolean;
	defaultTypeFilter: SupportedAttachment;
	defaultGroupFilter: string;
	defaultParentFilter: string;
	noCreateAttachment: boolean;
	displayOrderAttachment: number;
}

@Component({
	selector: 'splbk-list-media-dialog',
	templateUrl: './list-media.component.html',
	styleUrls: ['./list-media.component.scss']
})
export class ListMediaDialogComponent extends ListMediaComponent {

	constructor(spellbook: Spellbook, route: ActivatedRoute, screenService: ScreenService,
	            dialog: MatDialog, public dialogRef: MatDialogRef<ListMediaComponent>,
	            @Inject(MAT_DIALOG_DATA) private config: ListMediaDialogComponentConfig, private snackBar: MatSnackBar) {
		super(spellbook, route, screenService, dialog);

		this.multipleMode = false; // default
		if (config) {
			this.filterableByType = config.filterableByType;
			this.filterableByGroup = config.filterableByGroup;
			this.filterableByParent = config.filterableByParent;
			if (config.defaultTypeFilter) {
				this.filterType = {
					value: config.defaultTypeFilter.value,
					viewValue: config.defaultTypeFilter.name
				};
			}
			if (config.defaultGroupFilter) {
				this.filterGroup = {
					value: config.defaultGroupFilter,
					viewValue: config.defaultGroupFilter
				};
				this.group = config.defaultGroupFilter;
			}
			if (config.defaultParentFilter) {
				this.filterParent = {
					value: config.defaultParentFilter,
					viewValue: ListMediaComponent.MINE_LABEL
				};
				this.parentKey = config.defaultParentFilter;
			}
		}
	}

	showDetail(item: Attachment) {
		// Gallery: open delected item's detail
		const dialogRef = this.dialog.open(CreateUpdateMediaComponent, {
			width: '800px',
			data: item
		});
		dialogRef.componentInstance.noCreateAttachment = true;

		dialogRef.afterClosed().subscribe((a: Attachment) => {
			// come from content
			if (a !== undefined) {
				this.onItemClick(a);
			}
		});
	}

	onItemClick(item: Attachment, forceCreate?: boolean) {
		item.group = this.group;
		item.type = this.type;
		item.parentKey = this.parentKey;
		item.parentType = this.parentType;
		item.displayOrder = this.config.displayOrderAttachment;
		if (this.config && this.config.noCreateAttachment) {
			this.dialogRef.close(item);
		} else {
			// default: create also attachment
			// Create an attachment for this content
			this.isLoading = true;
			// console.log('onItemClick createAttachment', item);
			this.client.createAttachment(item).subscribe(
				(a: Attachment) => {
					this.isLoading = false;
					this.dialogRef.close(a);
				},
				(error) => {
					this.isLoading = false;
					this.snackBar.open(`Unable to create attachment`);
				}
			);
		}
	}

	showDetails(items: Attachment[]) {
		if (this.config && this.config.noCreateAttachment) {
			this.dialogRef.close(items);
			return;
		}
		// default: create also attachment
		// Create an attachment for this content
		const results = [];
		this.isLoading = true;
		items.forEach((item: Attachment, index: number) => {
			// const sameItem = this.isSameGroup(item) && this.isSameParent(item); // todo check same item
			item.group = this.group;
			item.type = this.type;
			item.parentKey = this.parentKey;
			item.parentType = this.parentType;
			item.displayOrder = this.config.displayOrderAttachment + index;
			results.push(this.client.createAttachment(item));

			if (results.length === items.length) {
				forkJoin(results).subscribe((atts: Attachment[]) => {
						this.isLoading = false;
						this.dialogRef.close(atts);
					},
					(error: HttpErrorResponse) => {
						console.error('error create attachment', error);
						this.isLoading = false;
						this.snackBar.open(`Unable to create attachment`);
					});
			}
		});
	}

}
