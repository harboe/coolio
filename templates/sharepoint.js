(function(win, doc, undefined) {

	{{ range $js := .Javascript }}
	{{$js}}
	{{ end }}

	var container = '<div id="coolio" class="container-fluid" data-bind="component: { name: \'group\', params: $data }"></div>';
	
	var postbox = new ko.subscribable();

	ko.subscribable.fn.subscribeTo = function(topic) {
		postbox.subscribe(this, null, topic);
		return this;  //support chaining
	};

	ko.subscribable.fn.publishOn = function(topic, broadcastOnInit) {
		this.subscribe(function(newValue) {
			postbox.notifySubscribers(newValue, topic);
		});

		return this; //support chaining
	};

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
					return 'radio'
				} else if (schema.FormatType === 2) {
					return 'checkbox'
				}
				return 'select'
			}
			if (t === 'lookup') {
				return 'select';
			}

			return t;
		};

		var registerFormContext = function(name, ctx) {
			var ctx = SPClientTemplates.Utility.GetFormContextForCurrentField(ctx);
			var schema = ctx.fieldSchema;
			
			spContext[name] = ctx;
			self[name] = {
				name: schema.Name,
				desc: schema.Description,
				type: formatFieldType(schema),
				aux: {
					readonly: schema.ReadOnlyField,
					required: schema.Required,
					maxLength: schema.MaxLength,
					rows: schema.NumberOfLines,
					options: getOptions(schema),
				},
				value: ctx.fieldValue
			};

			return ctx;
		};

		var templateOverride = {
			Templates: {
				Fields: {
					{{ range $param := .Parameters }}
					'{{$param.Id}}': {
						EditForm: function(ctx){
							registerFormContext('{{$param.Id}}', ctx);
						},
						NewForm: function(ctx){
							registerFormContext('{{$param.Id}}', ctx);
						}
					},
					{{ end }}
				}
			}
		};

		try { SPClientTemplates.TemplateManager.RegisterTemplateOverrides(templateOverride);}
		catch(err) { }
	};

	var DataModelDefaults = {
		id: '',
		help: '',
		params: [],
		groups: [],
		aux: {}
	};

	var DataModel = function(data) {
		$.extend(this, DataModelDefaults, data);

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
		data = $.extend({}, DataModelDefaults, coolio.sharepoint[data.id], data);

		this.id = ko.observable(data.id);
		this.name = ko.observable(data.name || data.id);
		this.desc = ko.observable(data.desc || data.name);
		this.help = ko.observable(data.help);
		this.type = ko.observable(data.type);
		this.value = ko.observable(data.value);
		this.validation = ko.observable('');
		this.aux = data.aux;

		if (data.id) {
			this.value.subscribeTo(data.id);
		}

		this.hasHelp = ko.computed(function() {
			return this.help().length > 0;
		}, this);

		this.params = ko.observableArray(data.params);
		this.groups = ko.observableArray(data.groups);

		coolio.sharepoint.register(data.id, this);
	};

	var Templates = function() {
		doc.write('{{.CustomHTML.Inline}}');
		{{ range $key, $val := .Templates }}
		doc.write('<template id="coolio-{{$key}}-template">{{$val.Inline}}</template>');
		ko.components.register('{{$key}}', {
			template: { element: 'coolio-{{$key}}-template' },
			{{ if $key.HasViewModel }} viewModel: ViewModel {{end}}
		});
		this['{{$key}}'] = true;
		{{ end}}
	};

	Templates.prototype.register = function(name, viewModel) {
		// if viewModel is undefined use ParameterViewMode;
		ko.components.register(name, {
			template: { element: name },
			viewModel: viewModel || ViewModel
		});
		this[name] = true;
	};

	var coolio = {
		templates: new Templates(),
		sharepoint: new Sharepoint(),
		data: new DataModel({{.JSON}}),
		subscribe: function(cb, context, topic) {
			postbox.subscribe(cb, context, topic);
		},
		notify: function(item, topic) {
			postbox.notifySubscribers(item, topic);
		}
	};

	try { {{.CustomJS}} }
	catch(err) { alert(err); }

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