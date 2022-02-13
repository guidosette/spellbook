import {HttpErrorResponse} from '@angular/common/http';
import {FormGroup} from '@angular/forms';
import {ResponseError} from './response-error';

/**
 * Created by Luigi Tanzini (luigi.tanzini@distudioapp.com) on 2019-07-11.
 */
export class ErrorUtils {

	static handlePostError(err: HttpErrorResponse, postForm?: FormGroup): ResponseError {
		let responseError: ResponseError;
		if (!err || !err.error) {
			responseError = new ResponseError(err.message, 'generic');
			return responseError;
		}
		if (err.error instanceof Array) {
			if (err.error.length === 0) {
				responseError = new ResponseError(err.message, 'generic');
				return responseError;
			}
			responseError = new ResponseError(err.error[0].error, err.error[0].field);
		} else {
			responseError = err.error;
		}
		if (postForm) {
			if (postForm.contains(responseError.Field)) {
				postForm.controls[responseError.Field].setErrors({incorrect: true});
			}
		}
		return responseError;
	}

}
