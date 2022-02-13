import {Component, ElementRef, NgZone, OnInit, ViewChild} from '@angular/core';
import {Spellbook} from '../core/spellbook';
import {ActivatedRoute} from '@angular/router';
import {MatDialog, MatSnackBar} from '@angular/material';
import {AbstractControl, FormControl, FormGroup, Validators} from '@angular/forms';

import {Place} from './place';

import {MapsAPILoader, MouseEvent} from '@agm/core';
import {HttpErrorResponse} from '@angular/common/http';
import {google} from '@google/maps';
import {PlaceClient} from './place-client';
import {ResponseError} from '../core/response-error';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {ErrorUtils} from '../core/error-utils';

declare var google: any;

@Component({
	selector: 'splbk-create-update-place',
	templateUrl: './create-update-place.component.html',
	styleUrls: ['./create-update-place.component.scss']
})
export class CreateUpdatePlaceComponent implements OnInit {

	private readonly client: PlaceClient;
	public place: Place;

	public get action(): string {
		return (!this.place.isNew()) ? 'Update' : 'Create';
	}

	public formGroup: FormGroup;
	public responseError: ResponseError;

	// MAP
	public zoom = 10;
	private geoCoder;
	@ViewChild('search')
	public searchElementRef: ElementRef;

	constructor(private spellbook: Spellbook, private route: ActivatedRoute, public dialog: MatDialog, private snackBar: MatSnackBar, private mapsAPILoader: MapsAPILoader,
	            private ngZone: NgZone) {
		this.client = new PlaceClient(spellbook);
		this.responseError = undefined;
		this.place = new Place();
	}

	ngOnInit() {
		this.buildForm();

		this.route.params.subscribe(params => {
			if (params.id) {
				this.place.id = Number(params.id);
				this.client.getPlace(this.place.id).subscribe(
					(p: Place) => {
						this.place = p;
						this.updateForm(p);
					},
					(error: HttpErrorResponse) => {
						console.error('Error', error);
						this.snackBar.open('Error! ' + error.statusText, 'ok', {});
					});
			} else {
				// create
				this.setCurrentLocation();
			}
		});

		this.formGroup.get('lat').valueChanges.subscribe(val => {
			this.place.lat = val;
			// this.getAddress(this.place.lat, this.place.lng);
		});
		this.formGroup.get('lng').valueChanges.subscribe(val => {
			this.place.lng = val;
			// this.getAddress(this.place.lat, this.place.lng);
		});

		this.initGeocode();

		this.formGroup.get('address').valueChanges.subscribe(val => {
			// this.geocodeService.geocodeAddress(val).subscribe((location: CustomLocation) => {
			// 	console.log('location', location);
			// });
		});

	}

	private buildForm(): void {
		this.formGroup = new FormGroup({
			name: new FormControl('', [Validators.required, Validators.minLength(2)]),
			address: new FormControl('', [Validators.required, Validators.minLength(2)]),
			street: new FormControl(''),
			streetNumber: new FormControl(''),
			area: new FormControl(''),
			city: new FormControl(''),
			postalCode: new FormControl(''),
			country: new FormControl(''),
			phone: new FormControl(''),
			description: new FormControl(''),
			lat: new FormControl('', [Validators.required]),
			lng: new FormControl('', [Validators.required]),
			website: new FormControl(''),
		});
	}

	private updateForm(place: Place) {
		this.formGroup.patchValue(place);
		// this.postForm.controls.order.setValue(1); // default
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.formGroup.controls[controlName].hasError(errorCode);
	}

	public hasMandatory(controlName: string): boolean {
		const formField = this.formGroup.get(controlName);
		if (!formField.validator) {
			return false;
		}
		const validator = formField.validator({} as AbstractControl);
		return (validator && validator.required);
	}

	public doCreateUpdate(formValue: any): void {
		if (this.formGroup.valid) {
			// populate the user object
			this.place.name = formValue.name;
			this.place.address = formValue.address;
			this.place.street = formValue.street;
			this.place.streetNumber = formValue.streetNumber;
			this.place.area = formValue.area;
			this.place.city = formValue.city;
			this.place.postalCode = formValue.postalCode;
			this.place.country = formValue.country;
			this.place.phone = formValue.phone;
			this.place.description = formValue.description;
			this.place.lat = formValue.lat;
			this.place.lng = formValue.lng;
			this.place.website = formValue.website;

			if (!this.place.isNew()) {
				this.update();
			} else {
				// create
				this.create();
			}
		}
	}

	private create(): void {

		this.responseError = undefined;
		this.client.createPlace(this.place).subscribe(
			(p: Place) => {
				this.snackBar.open('Place created!', 'ok', {});
				this.spellbook.router.navigate([`/place/${p.id}`]);
			},
			(err: HttpErrorResponse) => {
				this.place.id = 0;
				this.responseError = ErrorUtils.handlePostError(err, this.formGroup);
			}
		);
	}

	private update(): void {
		this.responseError = undefined;
		this.client.updatePlace(this.place).subscribe(
			(p: Place) => {
				this.place = p;
				this.updateForm(p);
				this.snackBar.open('Place updated!', 'ok', {});
			},
			(err: HttpErrorResponse) => {
				this.responseError = ErrorUtils.handlePostError(err, this.formGroup);
			}
		);
	}

	delete() {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + this.place.name + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			snackBarRef.dismiss();
			this.client.deletePlace(this.place).subscribe(
				() => {
					this.snackBar.open('Place deleted!', 'ok', {});
					// refresh
					this.spellbook.router.navigate(['/place']);
				},
				(err: HttpErrorResponse) => {
					this.responseError = ErrorUtils.handlePostError(err, this.formGroup);
				}
			);

		};
		snackbarData.actionNo = () => {
			snackBarRef.dismiss();
		};
	}

	// MAP
	markerDragEnd(place: Place, $event: MouseEvent) {
		console.log('dragEnd', place, $event);
		this.place.lat = $event.coords.lat;
		this.place.lng = $event.coords.lng;
		this.updateForm(this.place);
		this.getAddress(this.place.lat, this.place.lng);
	}

	clickedMarker(label: string) {
		console.log(`clicked the marker: ${label}`);
	}

	private setCurrentLocation() {
		if (window.navigator.geolocation) {
			window.navigator.geolocation.getCurrentPosition(
				position => {
					setTimeout(() => {
						this.getAddress(position.coords.latitude, position.coords.longitude);
					}, 0);
				},
				error => {
					switch (error.code) {
						case 1:
							console.log('Permission Denied');
							break;
						case 2:
							console.log('Position Unavailable');
							break;
						case 3:
							console.log('Timeout');
							break;
					}
				}
			);
		}
	}

	initGeocode() {
		// load Places Autocomplete
		this.mapsAPILoader.load().then(() => {
			this.geoCoder = new google.maps.Geocoder();

			const autocomplete = new google.maps.places.Autocomplete(this.searchElementRef.nativeElement, {
				types: ['address']
			});
			autocomplete.setFields([
				'formatted_address',
				'geometry',
				'address_components',
			]);
			autocomplete.addListener('place_changed', () => {
				this.ngZone.run(() => {
					// get the place result
					const place: google.maps.places.PlaceResult = autocomplete.getPlace();

					// verify result
					if (place.geometry === undefined || place.geometry === null) {
						return;
					}

					this.setPlace(place);
				});
			});
		});
	}

	setPlace(place: any): void {
		// set latitude, longitude and zoom
		this.place.lat = place.geometry.location.lat();
		this.place.lng = place.geometry.location.lng();
		this.place.address = place.formatted_address;
		this.setPlaceFromAddressComponents(place.address_components);
	}

	setPlaceFromAddressComponents(addresses: any) {
		this.place.city = this.getAddressComponent(addresses, 'locality');
		this.place.postalCode = this.getAddressComponent(addresses, 'postal_code');
		this.place.country = this.getAddressComponent(addresses, 'country');
		this.place.street = this.getAddressComponent(addresses, 'route');
		this.place.streetNumber = this.getAddressComponent(addresses, 'street_number');
		this.place.area = this.getAddressComponent(addresses, 'administrative_area_level_2');
		this.updateForm(this.place);
		this.zoom = 12;
	}

	getAddressComponent(addresses: any, type: string): string {
		let result = null;
		addresses.forEach((a: any) => {
			const longName = a.long_name;
			const types = a.types;
			const index = types.findIndex((t: string) => {
				return t === type;
			});
			if (index >= 0) {
				result = longName;
				return;
			}
		});
		return result;
	}

	getAddress(latitude, longitude) {
		if (!this.geoCoder) {
			return;
		}
		this.geoCoder.geocode({location: {lat: latitude, lng: longitude}}, (results, status) => {
			if (status === 'OK') {
				if (results[0]) {
					this.setPlace(results[0]);
				} else {
					console.error('No results found');
				}
			} else {
				console.error('Geocoder failed due to: ' + status);
			}

		});
	}
}
