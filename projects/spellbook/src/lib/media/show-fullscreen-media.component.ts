import {Component, OnInit} from '@angular/core';
import {Attachment} from './multimedia/attachment';
import {Spellbook} from '../core/spellbook';
import {ActivatedRoute} from '@angular/router';
import {MediaClient} from './media-client';
import {UrlUtils} from '../core/url-utils';
import {MatSnackBar} from '@angular/material/snack-bar';

@Component({
	selector: 'splbk-show-fullscreen-media',
	templateUrl: './show-fullscreen-media.component.html',
	styleUrls: ['./show-fullscreen-media.component.scss']
})
export class ShowFullscreenMediaComponent implements OnInit {
	private client: MediaClient;
	public item: Attachment;

	public utils = UrlUtils;


	constructor(private spellbook: Spellbook, private route: ActivatedRoute, private snackBar: MatSnackBar) {
		this.client = new MediaClient(spellbook);
	}

	ngOnInit() {
		this.route.params.subscribe(params => {
			if (params.id) {
				this.client.getAttachment(params.id).subscribe(
					(a: Attachment) => {
						this.item = a;
					},
					(error) => {
						console.error('Error', error);
						this.snackBar.open('Error! ' + error.statusText, 'ok', {});
					});
			}
		});
	}

	back() {
		this.spellbook.router.navigate([`../`], { relativeTo: this.route });
	}

}
