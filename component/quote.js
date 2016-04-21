/**
 * Created by yen-chieh on 4/13/16.
 */

var React = require('react');
var ReactDOM = require('react-dom');
var update = require('react-addons-update');
import 'whatwg-fetch';

const API_DOMAIN = "http://localhost:8080/";
const API_ADDSTOCK = "addToList";
const API_HEADER = {
	'Accept': 'application/json',
	'Content-Type': 'application/json'
};

var QuoteTable = React.createClass({
	getInitialState: function () {
		return {
			stocks: []
		}
	},

	componentWillReceiveProps: function(){
		if(this.props.data.length == 0){
			return;
		}
		var stocks = this.props.data;
		this.setState({stocks: stocks});
	},

	clickAddButton: function(e){
		var clickedButton = $(e.target);

		if(clickedButton.hasClass("disabled")){
			return;
		}

		clickedButton.addClass("disabled");
		this.storeStock({
			button: clickedButton
		});
	},

	storeStock: function(oArgs){
		oArgs = oArgs || {};
		var clickedButton = oArgs.button;
		var username = "Jay";
		var data = {
			symbol: clickedButton.attr('data-symbol'),
			stockName: clickedButton.attr('data-name'),
			username: username
		};

		fetch(API_DOMAIN + API_ADDSTOCK, {
			method: 'POST',
			mode: 'cors',
			body: JSON.stringify(data),
			headers: API_HEADER
		}).then((response) => response.json())
			.then((data) => {
				if(!data || data.length == 0){
					clickedButton.removeClass("disabled");

				}else{
					this.stockAdded(clickedButton);
				}
			}
		);
	},

	stockAdded: function(clickedButton){
		clickedButton.parents("tr").css({
			"background-color": "gainsboro"
		});

		let index = clickedButton.attr("data-id");
		let updateButton = update(this.state.stocks, {
			[index]: {
				buttonType: {
					$set: "remove"
				}
			}
		});

		this.setState({
			stocks: updateButton
		});

	},

	renderFirstTitle: function(){
		if(this.props.buttonType == "add"){
			return <th>Add</th>
		}else if(this.props.buttonType == "remove"){
			return <th>Remove</th>
		}else{
			return null
		}
	},

	renderAddOrRemove: function(quote, key){
		if(quote.buttonType == "add"){
			return (<td><button className="btn btn-sm btn-success" data-symbol={quote.symbol} data-name={quote.name} onClick={this.clickAddButton} data-id={key}>Add to List</button></td>)
		}else if(quote.buttonType == "remove"){
			return (<td><button className="btn btn-sm btn-warning" data-symbol={quote.symbol} data-name={quote.name} onClick={this.props.removeButtonCallback.bind(null, key)} data-id={key}>Remove</button></td>)
		}else{
			return null
		}
	},

	render: function () {
		if(!this.state.stocks){
			return;
		}

		var quoteTable = this.state.stocks.map(function (quote, i) {
				return (
					<tr key={i}>
						{this.renderAddOrRemove(quote, i)}
						<td>
							{quote.name}
						</td>
						<td>
							{quote.symbol}
						</td>
						<td>
							{quote.open}
						</td>
						<td>
							{quote.change} / {quote.percentChange}
						</td>
						<td>
							{quote.daysLow} / {quote.daysHigh}
						</td>
						<td>
							{quote.volume}
						</td>
					</tr>
				)

		}.bind(this));


		return (
			<div className="quote">
				<table className="table table-hover">
					<thead>
					<tr>
						{this.renderFirstTitle()}
						<th>Name</th>
						<th>Symbol</th>
						<th>Open</th>
						<th>Change / Percent</th>
						<th>DaysLow / DaysHigh</th>
						<th>Volume</th>
					</tr>
					</thead>

					<tbody>
					{quoteTable}
					</tbody>
				</table>

			</div>
		)
	}
});

module.exports = QuoteTable;