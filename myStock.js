/**
 * Created by yen-chieh on 4/16/16.
 */

var React = require('react');
var ReactDOM = require('react-dom');
var QuoteTable = require('./component/quote.js');

var MyStock = React.createClass({
	getInitialState: function () {
		return {
			userName: '',
			email: '',
			domain: 'http://localhost:8080/',
			getStockList: 'getUserStockList',
			stockList: {}
		}
	},

	componentDidMount: function () {
		var userName= "Jay";
/*		this.serverRequest = $.ajax({
			url: this.state.domain + this.state.getStockList,
			data: {
				userName: userName
			},
			method: 'GET',
			crossDomain: true,
			cache: false,
			dataType: 'json',
			success: function (data) {
				console.error(data);

			}.bind(this)
		});*/
	},

	componentWillMount: function () {
		//this.serverRequest.abort();
	},

	getUserStock: function () {

	},

	render: function () {
		return (
			<div>

			</div>
		)
	}
});

module.exports = MyStock;