{{ .JsLibraries }}
(function(win, doc, undefined) {

	{{ .Javascript }}

	var container = '<div id="coolio" class="container-fluid" data-bind="component: { name: \'group\', params: $data }"></div>';
	var templates = {{.Templates}};
	var overrides = {{.Overrides}};

	var Sharepoint = function() {
		var self = this;
		var spContext = {};

		this.register = function(name, viewModel) {
			var ctx = spContext[name];

			if (ctx) {
				ctx.registerGetValueCallback(ctx.fieldName, function() {
					return ko.unwrap(viewModel.value);
				});
				ctx.registerValidationErrorCallback(ctx.fieldName, function(a, b) {
					viewModel.validation(a.errorMessage);
				});
			}
		};

		var getOptions = function(schema) {
			var list = null;		

			if (schema.FieldType === 'Choice') {
				list = [];

				for (var i in schema.Choices) {
					var val = schema.Choices[i];
					list.push({ name: val, value: val});
				}
			} else if (schema.FieldType === 'Lookup') {
				list = [];

				for (var i in schema.Choices) {
					var val = schema.Choices[i];
					list.push({ name: val.LookupValue, value: val.LookupId });
				}
			}

			return list;
		};

		var formatFieldType = function(schema) {
			var t = schema.FieldType.toLowerCase();

			if (t === 'note')
				return 'textarea';
			if (t === 'choice') {
				if (schema.FormatType === 1) {
					return 'radio';
				} else if (schema.FormatType === 2) {
					return 'checkbox';
				}
				return 'select';
			}
			if (t === 'lookup') {
				return 'select';
			}

			return t;
		};

		var registerFormContext = function(name, ctx, readonly) {
			var fromCtx = SPClientTemplates.Utility.GetFormContextForCurrentField(ctx);

			if (!formCtx) {
				return
			}

			var schema = fromCtx.fieldSchema;
			var type = readonly ? 'static' : formatFieldType(schema);
			
			spContext[name] = fromCtx;
			self[name] = {
				name: schema.Name,
				desc: schema.Description,
				type: type,
				aux: {
					readonly: schema.ReadOnlyField,
					required: schema.Required,
					maxLength: schema.MaxLength,
					rows: schema.NumberOfLines,
					options: getOptions(schema)
				},
				value: fromCtx.fieldValue
			};
		};

		var fields = {};

		for(var i in overrides) {
			var id = overrides[i];
			fields[id] = {
				'EditForm': function(ctx) {
					registerFormContext(id, ctx, false);
				}, 
				'NewForm': function(ctx) {
					registerFormContext(id, ctx, false);
				},
				'DisplayForm': function (ctx) {
					registerFormContext(id, ctx, true);
				}
			};
		}

		try { SPClientTemplates.TemplateManager.RegisterTemplateOverrides({ Templates: { Fields: fields } });}
		catch(err) { }
	};

	var dataModelDefaults = {
		id: '',
		help: '',
		params: [],
		groups: [],
		aux: {}
	};

	var DataModel = function(data) {
		$.extend(this, dataModelDefaults, data);

		this.params = ko.utils.arrayMap(this.params, function(p) {
			return new DataModel(p);
		});
		this.groups = ko.utils.arrayMap(this.groups, function(p) {
			return new DataModel(p);
		});

		this.hasChildren = this.params.length > 0 || this.groups.length > 0;

		if (!this.type) {
			this.type = this.hasChildren ? 'group': 'text';
		}
	};

	var ViewModel = function(params) {
		var data = ko.toJS(params.data || params);
		data = $.extend({}, dataModelDefaults, coolio.sharepoint[data.id], data);

		this.id = ko.observable(data.id);
		this.name = ko.observable(data.name || data.id);
		this.desc = ko.observable(data.desc || data.name);
		this.help = ko.observable(data.help);
		this.type = ko.observable(data.type);
		this.value = ko.observable(data.value);
		this.validation = ko.observable('');
		this.aux = data.aux;

		if (data.id) {
			postbox.subscribe(function(newValue) {
				if (this.value() !== newValue) {
					this.value(newValue);
				}
			}, this, data.id);

			this.value.subscribe(function(newValue) {
				postbox.notifySubscribers(newValue, data.id);
			});
		}

		this.hasHelp = ko.computed(function() {
			return this.help().length > 0;
		}, this);

		this.params = ko.observableArray(data.params);
		this.groups = ko.observableArray(data.groups);

		coolio.sharepoint.register(data.id, this);
	};

	var Templates = function() {
		for (var i in templates) {
			var tmpl = templates[i];
			ko.components.register(tmpl.name, {
				template: tmpl.html,
				viewModel: tmpl.templateOnly ? null : ViewModel
			});
		}
	};

	Templates.prototype.register = function(name, viewModel) {
		// if viewModel is undefined use ParameterViewMode;
		ko.components.register(name, {
			template: { element: name },
			viewModel: viewModel || ViewModel
		});
	};

	var coolio = {
		templates: new Templates(),
		sharepoint: new Sharepoint(),
		data: new DataModel({{.JSON}}),
		subscribe: function(cb, context, topic) {
			postbox.subscribe(cb, context, topic);
		},
		notify: function(item, topic) {
			console.log('Postbox; Topic:' + topic + ', Value: ' + item);
			postbox.notifySubscribers(item, topic);
		}
	};

	{{if .CustomHTML}}doc.write('{{.CustomHTML.Inline}}');{{end}}
	{{if .CustomJS}}try { {{.CustomJS}} } catch(err) { alert(err); }{{end}}

	$().ready(function() {
		var b = $('body'), f = $('.ms-formtable');

		if (f.length === 0) {
			b.append(container);
		} else {
			$('#DeltaPlaceHolderMain').prepend(container);
		}

		console.log('coolio says elo!');
		console.log(coolio);
		ko.applyBindings(coolio, doc.getElementById('coolio'));
	});
})(window, document);