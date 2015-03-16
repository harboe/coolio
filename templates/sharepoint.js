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
		}
	},
	formContext: function(ctx) {
		try {
			return SPClientTemplates.Utility.GetFormContextForCurrentField(ctx);
		} catch(err) { console.log(err); }

		return {};
	},
	data: function(name) {
		var field = sharepoint.fields[name];

		if (field) {
			var schema = field.fieldSchema;

			return {
				name: schema.Name,
				desc: schema.Description,
				type: schema.FieldType.toLowerCase(),
				aux: {
					required: schema.Required,
					maxLength: schema.MaxLength
				},
				value: field.fieldValue
			}
		}

		return {};
	}
};

// after load.
sharepoint.register();