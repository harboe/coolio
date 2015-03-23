var postbox = new ko.subscribable();

ko.subscribable.fn.subscribeTo = function(topic) {
	postbox.subscribe(this, this, topic);
	return this;  //support chaining
};

ko.subscribable.fn.publishOn = function(topic) {
	this.subscribe(function(newValue) {
		postbox.notifySubscribers(newValue, topic);
	}, this);

	return this; //support chaining
};