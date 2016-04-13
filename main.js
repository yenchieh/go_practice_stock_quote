
var React = require('react');
var ReactDOM = require('react-dom');
var css = require("!style!css!sass!./scss/main.scss");


var Main = React.createClass({
	getInitialState: function () {
		return {
			domain: 'http://localhost:8080/',
			searchPath: 'search',
			$app: $('div#app')
		}
	},

	componentDidMount: function () {

	},

	searchSymbol: function(){
		var symbolInput = $('input#symbolSearchInput', this.state.$app);
		$.get(this.state.domain + this.state.searchPath + "?symbol=" + symbolInput.val(), function(data){
			console.error(data);
		})
	},

	render: function () {
		return (

			<div id="mainSearchComponent">
				<figure className="highlight">
					<div className="form-inline">
						<div className="form-group">
							<label for="symbolSearchInput" className="sr-only">Enter Stock Here</label>
							<input type="text" name="symbol" className="form-control" id="symbolSearchInput"
										 placeholder="Enter symbol or stock name"/>
						</div>
						<button className="btn btn-primary" onClick={this.searchSymbol}>Search</button>
					</div>


				</figure>
			</div>
		)
	}
});

ReactDOM.render(<Main/>, document.getElementById('app'));