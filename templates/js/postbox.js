var postbox = new ko.subscribable();

ko.observable.fn.subscribeTo = function(topic) {
	postbox.subscribe(this, this, topic);
	return this;  //support chaining
};

ko.observable.fn.publishOn = function(topic) {
	this.subscribe(function(newValue) {
		postbox.notifySubscribers(newValue, topic);
	}, this);

	return this; //support chaining
};