var gogeo = {
	defaults: {
		server: 'http://localhost:8080',
		provider: 'google',
		baseUrl: function(format) {
			var url = this.server + '/v1/' + this.provider + '/' + format;

			if (this.apikey) {
				url += '?key=' = this.apikey;
			}

			return url;
		},
		addressUrl: function(format, address) {
			url = this.baseUrl(format);
			url += '?addr=' + address;

			return encodeURI(url);
		},
		locationUrl: function(format, location, args) {
			url = this.baseUrl(format);
			url += '?loc=' + location.lat + ',' + location.lng;

			for (var i in args) {
				url += '&' + i + '=' + args[i];
			}

			return encodeURI(url);
		}
	},
	lookup: function(url, callback, ctx) {
		var responseHandler = function(e, status, xhr) {
			if (status === 'error') {
				callback.apply(ctx, [null, e, xhr]);
			} else {
				callback.apply(ctx, [e, null, xhr]);
			}
		};

		$.ajax({
			method: 'GET',
			url: url,
			success: responseHandler,
			error: responseHandler
		});
	},
	init: function(elm, valueAccessor, allBindings, viewModel, ctx) {
		var $elm = $(elm);
		var val = valueAccessor();

		var opts = allBindings.get('gogeoOptions');
		opts = ko.utils.extend(opts, gogeo.defaults);

		if (!val) {
			throw new Error('missing value bindings');
		}

		var query = ko.observable();
		var newValueAccessor = function() {
			return query;
		};
		var self = false;
		var lookupAddress = function(newAddress) {
			if (self) {
				return
			}

			if (!newAddress || newAddress === '') {
				val(null);
				query('');
			} else {
				self = true;
				var url = opts.addressUrl('json', newAddress);
				gogeo.lookup(url, function(e, err) {
					if (err) {
						val(null);
						query('');

						console.log(err);
						return;
					}
					var geo = JSON.parse(e);
					var addr = geo && geo.length > 0 ? geo[0] : {};

					val(addr);
					query(addr.address);
					self = false;
				}, this);
			}
		};
		var image = function(newAddress) {
			// remove the image.
			$elm.next().remove();

			if (newAddress) {
				var url = opts.locationUrl('png', newAddress.location);
				$elm.parent().append('<br><img src="' + url + '" class="img-thumbnail img-responsive center-block" />');
			}
		};

		query.subscribe(lookupAddress);
		val.subscribe(lookupAddress);

		return ko.bindingHandlers.value.init(elm, newValueAccessor, allBindings, viewModel, ctx);
	}
};

ko.bindingHandlers['gogeo'] = gogeo;