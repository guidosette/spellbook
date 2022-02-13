import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef, MatSnackBar} from '@angular/material';
import {Observable} from 'rxjs';
import {Attachment} from './multimedia/attachment';
import {MediaClient} from './media-client';
import {Spellbook} from '../core/spellbook';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {UrlUtils} from '../core/url-utils';

@Component({
	selector: 'splbk-create-update-media',
	templateUrl: './create-update-media.component.html',
	styleUrls: ['./create-update-media.component.scss']
})
export class CreateUpdateMediaComponent implements OnInit {
	private readonly client: MediaClient;
	public item: Attachment;

	public utils = UrlUtils;

	public showImageFullScreen = false;
	public noCreateAttachment = false;

	constructor(private spellbook: Spellbook, @Inject(MAT_DIALOG_DATA) public attachment: Attachment,
	            public dialogRef: MatDialogRef<CreateUpdateMediaComponent>, private snackBar: MatSnackBar) {
		this.client = new MediaClient(spellbook);
		this.item = attachment;
	}

	ngOnInit() {}

	public save() {
		// console.log('save', this.item);

		if (!this.item.parentKey) {
			this.item.parentKey = Attachment.PARENT_DEFAULT;
		}
		if (!this.item.type) {
			this.item.type = Attachment.TYPE_DEFAULT;
		}

		if (!this.item.group) {
			this.item.group = Attachment.GROUP_DEFAULT;
		}

		if (this.noCreateAttachment) {
			// no create
			this.dialogRef.close(this.item);
			return;
		}
		let obs: Observable<Attachment>;
		if (this.item.id) {
			obs = this.client.updateAttachment(this.item);
		} else {
			obs = this.client.createAttachment(this.item);
		}
		obs.subscribe(
			(a: Attachment) => {
				this.snackBar.open('Media saved!', 'ok', {});
				this.dialogRef.close(a);
			},
			(err) => {
				this.snackBar.open('Error saving media', 'ok', {});
				console.log(err.Error);
			});
	}

	public delete() {
		if (this.item.id) {
			const snackbarData: SnackbarData = new SnackbarData();
			snackbarData.message = 'Are you sure to delete ' + this.item.name + '?';
			const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
				duration: 30000,
				data: snackbarData
			});
			snackbarData.actionOk = () => {
				snackBarRef.dismiss();
				this.client.deleteAttachment(this.item).subscribe(
					() => {
						this.snackBar.open('Media deleted!', 'ok', {});
						this.item.id = undefined;
						this.dialogRef.close(this.item);
					},
					(err) => {
						this.snackBar.open('Error deleting media', 'ok', {});
						console.log(err.Error);
					});
			};
			snackbarData.actionNo = () => {
				snackBarRef.dismiss();
			};
		}
	}

	public closeDialog() {
		this.dialogRef.close();
	}

	public showImageFullscreen() {
		if (this.showImageFullScreen && this.item.id !== undefined && this.spellbook.router) {
			this.spellbook.router.navigate([`media/${this.item.id}`]);
			this.closeDialog();
		}
	}

	public getIconForAttachmentName(name: string) {
		return Attachment.getIconForAttachmentName(name);
	}
}
