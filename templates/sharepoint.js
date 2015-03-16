// sharepoint blackmagic
var sharepoint = {
	fields: {},
	Templates: {
		Fields: {
			{{ range $param := .Parameters }}
			'{{$param.Id}}': {
				EditForm: function(ctx){
					sharepoint.fields['{{$param.Id}}'] = sharepoint.formContext(ctx);
				},
				NewForm: function(ctx){
					sharepoint.fields['{{$param.Id}}'] = sharepoint.formContext(ctx);
				}
			},
			{{ end }}
		}
	}, 
	register: function() {
		try { SPClientTemplates.TemplateManager.RegisterTemplateOverrides(sharepoint);}
		catch(err) { console.log(err); }
	},
	registerValue: function(name, viewModel) {
		var ctx = sharepoint.fields[name];

		if (ctx) {
			ctx.registerGetValueCallback(ctx.fieldName, function() {
				return ko.unwrap(viewModel.value);
			});
			ctx.registerValidationErrorCallback(ctx.fieldName, function(a, b) {
				viewModel.validation(a.errorMessage);
			});
		}
	},
	formContext: function(ctx) {
		try {
			return SPClientTemplates.Utility.GetFormContextForCurrentField(ctx);
		} catch(err) { console.log(err); }

		return {};
	},
	optionList: function(schema) {
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
	},
	typeName: function(schema) {
		var t = schema.FieldType.toLowerCase();

		switch (t) {
			case 'note':
			return 'textarea'
			case 'choice':
			if (schema.FormatType === 1) {
				return 'radio'
			} else if (schema.FormatType === 2) {
				return 'checkbox'
			}
			return 'select'
			case 'lookup':
			return 'select'
			default:
			return t
		}
	},
	data: function(name) {
		var field = sharepoint.fields[name]

		if (field) {
			var schema = field.fieldSchema;

			return {
				name: schema.Name,
				desc: schema.Description,
				type: this.typeName(schema),
				aux: {
					readonly: schema.ReadOnlyField,
					required: schema.Required,
					maxLength: schema.MaxLength,
					rows: schema.NumberOfLines,
					options: this.optionList(schema),

				},
				value: field.fieldValue
			}
		}

		return {};
	}
};

// after load.
sharepoint.register();