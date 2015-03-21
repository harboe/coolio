(function() {
	var resize = function() {
		var height = $('body').height() - $('nav').outerHeight() - 70;
		$('.fill, .CodeMirror').css('height', height + 'px');
	};
	var throttle_timer, yamlCM, htmlCM, jsCM, iframe, viewname;
	var save = function(preview) {
		clearTimeout(throttle_timer);
		throttle_timer = setTimeout(function() {
			var data = { yaml: yamlCM.getValue(), html: htmlCM.getValue(), js: jsCM.getValue() };

			$.post('/' + viewname + (preview === true ? '?preview' : ''), data, function(e) {
				iframe.srcdoc = e;
			});
		}, 500);

		return false;
	};
	var run = function() {
		save(true);
		return false;
	};
	var showYaml = function() {
		$('#html, #js').hide();
		$('#yaml').show();

		yamlCM.refresh();
		yamlCM.focus();
		return false;
	};

	var showHtml = function() {
		$('#js, #yaml').hide();
		$('#html').show();

		htmlCM.refresh();
		htmlCM.focus();
		return false;
	};

	var showJs = function() {
		$('#html, #yaml').hide();
		$('#js').show();
		$('a[href="#js"]').addClass('')

		jsCM.refresh();
		jsCM.focus();
		return false;
	};
	var showResult = function() {
		return false;
	};

	$().ready(function() {
		iframe = document.getElementById('result');
		viewname = $('#viewname').val();

		$('a[href="#save"]').click(save);
		$('a[href="#run"]').click(run);

		$('a[href="#result"]').click(function(e) {
			$('#yaml, #html, #js').hide();
			$('#result').show();

			e.preventDefault();
			return false;
		});
		$('a[href="#yaml"], a[href="#html"], a[href="#js"]').click(function(e) {
			$('#yaml, #html, #js').hide();
			var elm = $(this.hash)
			elm.show();

			var cm = elm.find('.CodeMirror')[0].CodeMirror;
			cm.refresh();
			cm.focus();

			e.preventDefault();
			return false;
		});

		yamlCM = CodeMirror.fromTextArea($('#yaml textarea')[0], {
			lineNumbers: true,
			extraKeys: { Tab: function(cm) {
				if (cm.somethingSelected()) {
					cm.indentSelection("add");
				} else {
					cm.replaceSelection(cm.getOption("indentWithTabs")? "\t":
						Array(cm.getOption("indentUnit") + 1).join(" "), "end", "+input");
				}
			}}
		});
		yamlCM.on('change', run);
		jsCM = CodeMirror.fromTextArea($('#js textarea')[0], {
			lineNumbers: true
		});

		htmlCM = CodeMirror.fromTextArea($('#html textarea')[0], {
			lineNumbers: true
		});

		// init preview
		run();

		$('.CodeMirror').addClass('border');

		resize();
		$(window).resize(resize);
	});
})();