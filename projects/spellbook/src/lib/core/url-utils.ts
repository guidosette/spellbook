/**
 * Created by Luigi Tanzini (luigi.tanzini@distudioapp.com) on 2019-07-11.
 */
export class UrlUtils {

	static formatBackgroundUrl(url: string, encode: boolean = false) {
		const encoded = encodeURI(url);
		return `url('${encode ? encoded : url}')`;
	}

}
