/**
 * Created by yen-chieh on 4/13/16.
 */

var React = require('react');
var ReactDOM = require('react-dom');

var QuoteTable = React.createClass({
	render: function () {
		if(!this.props.data){
			return;
		}
		var quoteTable = this.props.data.map(function (quotes, i) {
			return quotes.map(function(quote, i){
				return (
					<tr key={i}>
						<td>
							<button className="btn btn-sm btn-success">Add to List</button>
						</td>
						<td>
							{quote.resource.fields.name}
						</td>
						<td>
							{quote.resource.fields.price}
						</td>
						<td>
							{quote.resource.fields.symbol}
						</td>
						<td>
							{quote.resource.fields.volume}
						</td>
					</tr>
				)
			});

		});


		return (
			<div className="quote">
				<table className="table table-hover">
					<thead>
					<tr>
						<th>Add</th>
						<th>Name</th>
						<th>Price</th>
						<th>Symbol</th>
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