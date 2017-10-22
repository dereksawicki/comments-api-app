import React from 'react'
import Comment from './Comment'

class CommentsList extends React.Component {
	constructor(props) {
		super(props);

	}

	render() {
		return (
			<div className="comments-list">
				<h2>Comments</h2>
				<ul>
					{this.props.comments.map( (val) =>
					<li key={val['id']}> 
						<Comment id={val['id']} 
								 username={val['username']} 
								 body={val['body']} 
								 modtime={val['mod_time']}
								 onDelete={this.props.onDelete}
								 onEdit={this.props.onEdit}/>
					</li>
					)}
				</ul>
			</div>
		);
	}
	
}

export default CommentsList