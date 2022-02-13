import {Spellbook} from './spellbook';
import {FileAttachment} from '../media/multimedia/file';
import {HttpErrorResponse} from '@angular/common/http';

export class CkUploadAdapter {
	private loader;
	public loadingImage: boolean;
	public errorUploading: string;

	constructor(private spellbook: Spellbook, private type: string, private namespace: string, loader: any) {
		this.loader = loader;
	}

	/*
	Interface UploadAdapter
	 */
	public upload(): Promise<any> {
		let resolvePromise;
		let rejectPromise;
		this.loadingImage = true;
		this.loader.file.then(
			(data) => {
				console.log('data', data);
				const file = data;
				this.errorUploading = '';
				this.spellbook.api.postFile(
					this.type, this.namespace, file.name, file
				).subscribe((res: FileAttachment) => {
					const obj = {
						default: res.resourceUrl
					};
					this.loadingImage = false;
					resolvePromise(obj);
				}, (error: HttpErrorResponse) => {
					console.error('Error upload', error);
					this.loadingImage = false;
					this.errorUploading = error.statusText;
					rejectPromise(error);
				});
			});
		return new Promise((resolve, reject) => {
			resolvePromise = resolve;
			rejectPromise = reject;
		});
	}
}
