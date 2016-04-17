/**
 * Created by yen-chieh on 4/13/16.
 */

var React = require('react');
var ReactDOM = require('react-dom');

var QuoteTable = React.createClass({
	getInitialState: function () {
		return {
			domain: 'http://localhost:8080/',
			addStock: 'addToList'
		}
	},

	clickAddButton: function(e){
		this.clickedButton = $(e.target);
		var username = "Jay";
		this.clickedButton.addClass("disabled");
		this.storeStock({
			symbol: this.clickedButton.attr('data-symbol'),
			stockName: this.clickedButton.attr('data-name'),
			username: "Jay"
		});
	},

	storeStock: function(data){

/*
		this.get(this.state.domain + this.state.addStock, data).then(function(result){
			console.error(result);
		}, function(e){
			console.error(e);
		});
*/

		$.ajax({
			url: this.state.domain + this.state.addStock,
			data: data,
			method: 'GET',
			crossDomain: true,
			cache: false,
			dataType: 'json',
			success: function(data){

				if(!data || data.length == 0){
					this.clickedButton.removeClass("disabled");
				}else{
					this.clickedButton.parents("tr").css({
						"background-color": "gainsboro"
					});
				}
			}.bind(this)
		});
	},

	get: function(url, data){
		return new Promise(function(resolve, reject){
			var req = new XMLHttpRequest();
			req.open('GET', url, true);
			req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
			req.setRequestHeader("Access-Control-Allow-Origin", "*");
			req.setRequestHeader("Access-Control-Allow-Credentials", "true");
			req.setRequestHeader("Access-Control-Allow-Methods", "GET");
			req.setRequestHeader('Cache-Control', 'no-cache');
			req.onload = function() {
				if (req.status == 200) {
					resolve(req.response);
				} else {
					reject(Error(req.statusText));
				}
			};

			req.onerror = function(){
				reject(Error("Network Error"));
			};

			req.send(JSON.stringify(data));

		}.bind(this));
	},

	render: function () {
		if(!this.props.data){
			return;
		}
		var quoteTable = this.props.data.map(function (quote, i) {
				return (
					<tr key={i}>
						<td>
							<button className="btn btn-sm btn-success" data-symbol={quote.symbol} data-name={quote.name} onClick={this.clickAddButton}>Add to List</button>
						</td>
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
						<th>Add</th>
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