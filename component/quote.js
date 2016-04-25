/**
 * Created by yen-chieh on 4/13/16.
 */

var React = require('react');
var ReactDOM = require('react-dom');
var update = require('react-addons-update');
var ReactCSSTransitionGroup = require('react-addons-css-transition-group');
import 'whatwg-fetch';


var QuoteTable = React.createClass({
	getInitialState: function () {
		return {
			stocks: []
		}
	},

	componentWillReceiveProps: function(props){
		this.setState({stocks: props.data});
	},

	renderAddOrRemove: function(quote, key){
		if(quote.buttonType == "add"){
			return (<td><li className="btn btn-sm btn-success" onClick={this.props.addButtonCallback.bind(null, quote)} data-id={key}>Add to List</li></td>)
		}else if(quote.buttonType == "remove"){
			return (<td><button className="btn btn-sm btn-warning" onClick={this.props.removeButtonCallback.bind(null, key)}>Remove</button></td>)
		}else{
			return null
		}
	},

	renderColoredStockChange: function(change, percentChange){
		let className = change < 0 ? "negative" : "positive";
		return(
			<div>
				<span className={className}>{change}</span> / <span className={className}>{percentChange}</span>
			</div>

		)
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
							<div className="symbol">{quote.symbol}</div>
						</td>
						<td>
							{quote.open}
						</td>
						<td>
							{this.renderColoredStockChange(quote.change, quote.percentChange)}
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
						<th>Options</th>
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

export default QuoteTable;