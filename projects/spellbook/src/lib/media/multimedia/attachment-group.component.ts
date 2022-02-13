import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Spellbook} from '../../core/spellbook';
import {MatDialog, MatSnackBar} from '@angular/material';
import {AttachmentGroup} from './attachment-group';
import {ScreenService} from '../../core/screen.service';

import {Attachment} from './attachment';
import {MediaClient} from '../media-client';
import {SupportedAttachment} from '../../core/supported-attachment';
import {ListMediaDialogComponent} from '../list-media-dialog.component';
import {CreateUpdateMediaComponent} from '../create-update-media.component';
import {SnackbarComponent, SnackbarData} from '../../core/snackbar.component';
import {UrlUtils} from '../../core/url-utils';
import {moveItemInArray} from '@angular/cdk/drag-drop';
import {HttpErrorResponse} from '@angular/common/http';
import {ErrorUtils} from '../../core/error-utils';
import {Subscription} from 'rxjs';
import {DragulaService} from 'ng2-dragula';

@Component({
	selector: 'splbk-multimedia-group',
	templateUrl: './attachment-group.component.html',
	styleUrls: ['./attachment-group.component.css']
})
export class AttachmentGroupComponent implements OnInit, OnDestroy {
	private readonly client: MediaClient;

	@Input()
	public attachmentGroup: AttachmentGroup;

	@Input()
	public max: number; // if not set = infinite

	@Input()
	public type: string;

	public supportedTypes: Array<SupportedAttachment>;

	panelOpenState = false;

	public utils = UrlUtils;

	// dragula
	subs = new Subscription();

	public readonly group: string;

	constructor(private spellbook: Spellbook, public dialog: MatDialog, private screenService: ScreenService, private snackBar: MatSnackBar, private dragulaService: DragulaService) {
		this.client = new MediaClient(spellbook);
		this.supportedTypes = this.spellbook.supportedAttachments;

		// this.screenService.screenSize.asObservable().subscribe((screenSize: ScreenSize) => {
		// 	this.screenWidth = screenSize.width;
		// 	// todo
		// });

		this.group = spellbook.uniqueId();
		this.initDragula();
	}


	ngOnInit() {
	}

	canAdd(): boolean {
		if (this.max > 0 && this.attachmentGroup.attachments.length >= this.max) {
			return false;
		}
		return true;
	}


	getImageForType(type: string) {
		const support = this.supportedTypes.find((s) => {
			return s.value === type;
		});
		if (support) {
			return support.image;
		}
		return null;
	}

	addAttachment() {
		let defaultTypeFilter: SupportedAttachment;
		for (const sa of this.spellbook.supportedAttachments) {
			if (sa.value === this.attachmentGroup.type) {
				defaultTypeFilter = sa;
			}
		}
		const dialogRef = this.dialog.open(ListMediaDialogComponent, {
			width: '90vw',
			height: '90vh',
			data: {
				filterableByType: false,
				filterableByGroup: true,
				filterableByParent: true,
				defaultGroupFilter: this.attachmentGroup.name,
				defaultParentFilter: this.attachmentGroup.parentKey.toString(),
				noCreateAttachment: false, // create
				displayOrderAttachment: this.attachmentGroup.attachments.length,
				defaultTypeFilter,
			}
		});
		dialogRef.componentInstance.group = this.attachmentGroup.name;
		dialogRef.componentInstance.type = this.attachmentGroup.type;
		dialogRef.componentInstance.parentKey = `${this.attachmentGroup.parentKey}`;
		dialogRef.componentInstance.type = `${this.attachmentGroup.type}`;
		dialogRef.componentInstance.parentType = this.type;
		dialogRef.componentInstance.multipleMode = true;
		dialogRef.afterClosed().subscribe((results: Attachment[]) => {
			// results: Attachment[] o Attachment
			console.log('addAttachment', results);
			if (results) {
				if (results instanceof Array) {
					// multiple
					results.forEach((result: Attachment) => {
						this.attachmentGroup.addAttachment(result);
					});
					this.snackBar.open(results.length + ' attachments added', 'ok', {});
					// for safety: all
					// this.refreshAllOrders();
				} else {
					// single
					const result: Attachment = results;
					this.attachmentGroup.addAttachment(result);
					this.snackBar.open('Attachment added', 'ok', {});
					// for safety: all
					// this.refreshAllOrders();
				}
			}
		});
	}

	onItemClick(attachment: Attachment) {
		const dialogRef = this.dialog.open(CreateUpdateMediaComponent, {
			width: '900px',
			data: attachment
		});

		dialogRef.afterClosed().subscribe((result: Attachment) => {
			if (result !== undefined) {
				if (result.id) {
					// Item updated
					const index = this.getAttachmentIndex(result.id);
					if (index >= 0) {
						this.attachmentGroup.attachments[index] = result;
					}
				} else {
					// Item deleted
					const index = this.getAttachmentIndex(attachment.id);
					if (index >= 0) {
						this.attachmentGroup.attachments.splice(index, 1);
					}
					this.updateAttachmentsOrder(index, (this.attachmentGroup.attachments.length - 1));
				}
			}
		});

	}

	deleteAttachment(attachment: Attachment) {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + attachment.name + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			snackBarRef.dismiss();
			this.client.deleteAttachment(attachment).subscribe(
				() => {
					const index = this.getAttachmentIndex(attachment.id);
					if (index >= 0) {
						this.attachmentGroup.attachments.splice(index, 1);
					}
					this.updateAttachmentsOrder(index, (this.attachmentGroup.attachments.length - 1));
				},
				(error) => {
					this.snackBar.open(`Unable to delete attachment ${attachment.name}`);
				}
			);
		};
		snackbarData.actionNo = () => {
			console.log('actionNo');
			snackBarRef.dismiss();
		};
	}

	private getAttachmentIndex(id: string) {
		return this.attachmentGroup.attachments.findIndex((att: Attachment) => {
			return att.id === id;
		});
	}

	public getIconForAttachmentName(name: string) {
		return Attachment.getIconForAttachmentName(name);
	}

	// dragula

	initDragula() {
		this.dragulaService.createGroup(this.group, {
			moves: (el, container, handle) => {
				return handle.classList.contains('handle');
			}
		});

		this.subs.add(this.dragulaService.dropModel(this.group)
			.subscribe(({el, target, source, sourceModel, targetModel, item}) => {
				// console.log('***dropModel:');
				// console.log('el', el);
				const att: Attachment = item;
				const arr: Array<Element> = Array.prototype.slice.call(source.children);
				const idx = arr.indexOf(el);
				if (idx === -1) {
					return;
				}
				const from = att.displayOrder || this.attachmentGroup.attachments.indexOf(att);
				this.drop(from, idx);
			})
		);
		this.subs.add(this.dragulaService.removeModel(this.group)
			.subscribe(({el, source, item, sourceModel}) => {
			})
		);
	}

	drop(previousIndex: number, currentIndex: number) {
		// console.log('drop', previousIndex, currentIndex);
		moveItemInArray(this.attachmentGroup.attachments, previousIndex, currentIndex);
		const minIndex = Math.min(previousIndex, currentIndex);
		const maxIndex = Math.max(previousIndex, currentIndex);
		this.updateAttachmentsOrder(minIndex, maxIndex);
	}

	refreshAllOrders() {
		this.updateAttachmentsOrder(0, (this.attachmentGroup.attachments.length - 1));
	}

	updateAttachmentsOrder(min: number, max: number) {
		// console.log('updateAttachmentsOrder', min, max);
		const minIndex = Math.min(min, 0);
		const maxIndex = Math.max(max, this.attachmentGroup.attachments.length - 1);

		this.attachmentGroup.attachments.forEach((att: Attachment, index: number) => {
			if (index >= minIndex && index <= maxIndex) {
				console.log('att', index, att.displayOrder);
				att.displayOrder = index;
				this.client.updateAttachment(att).subscribe(
					(m: Attachment) => {
					},
					(err: HttpErrorResponse) => {
						const responseError = ErrorUtils.handlePostError(err);
						console.error('error updateAttachment', responseError.Error);
						this.snackBar.open(responseError.Error, 'ok', {});
					}
				);
			}
		});
	}

	ngOnDestroy() {
		this.subs.unsubscribe();
		this.dragulaService.destroy(this.group);
	}

}
