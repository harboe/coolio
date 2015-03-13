ko.bindingHandlers.toggle = {
	init: function(elm) {
		var $elm = $(elm), toggle = $elm.data('toggle'), hash = elm.hash;
		// console.log('init bootstrap: ' + toggle);

		if (toggle === 'collapse') {
			$elm.click(function(e) {
				e.preventDefault();
				$(hash).collapse({toggle: true});
			});
		} else if (toggle === 'tab') {
			$elm.click(function(e) {
				$elm.tab('show');

				// fix issue with when triggering data-toggle
				// outside of the nav-tabs.
				var tab = $('.nav-tabs a[href="' + hash + '"]').parent();
				tab.parent().find('.active').removeClass('active');
				tab.addClass('active');

				e.preventDefault();
			});
		} else if (toggle === 'tooltip') {
			$elm.tooltip();
			console.log('hay...!');
		}
	}
};