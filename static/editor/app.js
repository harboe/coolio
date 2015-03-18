var resize = function() {
	var height = $('body').height() - $('nav').outerHeight() - 40;
	$('.fill, .CodeMirror').css('height', height + 'px');
};
var throttle_timer, yamlCM, htmlCM, jsCM, iframe, viewname;
var saveData = function(preview) {
	clearTimeout(throttle_timer);
	throttle_timer = setTimeout(function() {
		var data = { yaml: yamlCM.getValue(), html: htmlCM.getValue(), js: jsCM.getValue() };

		$.post('/' + viewname + (preview ? '?preview' : ''), data, function(e) {
			iframe.srcdoc = e;
		});
	}, 500);
};

$().ready(function() {
	iframe = document.getElementById('result');
	viewname = $('#viewname').val();

	$('a[href="#save"]').click(function(e) {
		saveData();
		e.preventDefault();
		return false;
	});
	$('a[href="#run"]').click(function(e) {
		saveData(true);
		e.preventDefault();
		return false;
	});

	$('a[href="#result"]').click(function(e) {
		$('#yaml, #html, #js').hide();
		$('#result').show();

		e.preventDefault();
		return false;
	});
	$('a[href="#yaml"], a[href="#html"], a[href="#js"]').click(function(e) {
		$('#yaml, #html, #js').hide();
		$(this.hash).show();

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
	yamlCM.on('change', function() { saveData(true); });
	jsCM = CodeMirror.fromTextArea($('#js textarea')[0], {
		lineNumbers: true
	});

	htmlCM = CodeMirror.fromTextArea($('#html textarea')[0], {
		lineNumbers: true
	});

	saveData(true);

	$('.CodeMirror').addClass('border');

	resize();
	$(window).resize(resize);
});