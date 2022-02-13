import {AfterViewInit, Component, Inject, OnInit, ViewChild} from '@angular/core';
import {FileAttachment} from './multimedia/file';
import {MAT_DIALOG_DATA, MatDialogRef, MatPaginator, MatSnackBar} from '@angular/material';
import {Attachment} from './multimedia/attachment';
import {SupportedAttachment} from '../core/supported-attachment';
import {Spellbook} from '../core/spellbook';
import {UrlUtils} from '../core/url-utils';
import {ListResponse} from '../core/client';
import {of as observableOf} from 'rxjs';
import {catchError, map, startWith, switchMap} from 'rxjs/operators';

@Component({
	selector: 'splbk-browse-media-files',
	templateUrl: './browse-media-files.component.html',
	styleUrls: ['./browse-media-files.component.scss']
})
export class BrowseMediaFilesComponent implements OnInit, AfterViewInit {

	constructor(private spellbook: Spellbook, private snackBar: MatSnackBar,
	            public dialogRef: MatDialogRef<BrowseMediaFilesComponent>,
	            @Inject(MAT_DIALOG_DATA) private attachment: Attachment) {
		for (const sa of this.spellbook.supportedAttachments) {
			if (sa.value === attachment.type) {
				this.selectedType = sa;
				break;
			}
		}
	}

	urlFiles: FileAttachment[] = [];
	public isLoading: boolean;
	readonly selectedType: SupportedAttachment;

	public utils = UrlUtils;

	public currentSize: number;
	@ViewChild(MatPaginator) paginator: MatPaginator;

	ngOnInit() {
	}

	ngAfterViewInit() {
		setTimeout(() => {
			this.load(0);
		});
	}

	public load(index: number) {
		this.paginator.pageIndex = index;
		this.paginator.page
			.pipe(
				startWith({}),
				switchMap(() => {
					this.isLoading = true;
					return this.spellbook.api.getUrlFiles(this.paginator.pageIndex, this.paginator.pageSize);
				}),
				map((response: ListResponse<FileAttachment>) => {
					this.isLoading = false;
					const items = response.items;
					const size = (this.paginator.pageIndex + 1) * Math.min(this.paginator.pageSize, items.length);
					this.currentSize = response.more ? size + 1 : size;
					return items;
				}),
				catchError(() => {
					this.isLoading = false;
					this.snackBar.open('Error: unable to retrieve files from server');
					return observableOf([]);
				})
			).subscribe(
			(fileAttachments: FileAttachment[]) => {
				fileAttachments.forEach((f: FileAttachment) => {
					// update attachmentType
					f.attachmentType = this.spellbook.getAttachmentTypeForFileAttachment(f);
				});
				this.urlFiles = fileAttachments;
			});
	}

	public select(item: FileAttachment) {
		if (this.isFileAllowed(item)) {
			item.attachmentType = this.spellbook.getAttachmentTypeForFileAttachment(item);
			this.dialogRef.close(item);
		}
	}

	public isFileAllowed(item: FileAttachment) {
		if (!this.selectedType || this.selectedType.accept === '*/*') {
			return true;
		}
		const attachmentType = this.spellbook.getAttachmentTypeForFileAttachment(item);
		return attachmentType === this.selectedType.value;
	}

	public getIconForAttachmentName(name: string) {
		return Attachment.getIconForAttachmentName(name);
	}

}
