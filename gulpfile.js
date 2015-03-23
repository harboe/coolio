var del = require('del')
, gulp = require('gulp')
, concat = require('gulp-concat')
, htmlreplace = require('gulp-html-replace')
, gulpif = require('gulp-if')
, minifyCss = require('gulp-minify-css')
, replacePath = require('gulp-replace-path')
, uglify = require('gulp-uglify')
, useref = require('gulp-useref')
, util = require('gulp-util');

var env = (util.env.env || 'development');
var dir = 'dist';

// Copies index.html, replacing <script> and <link> tags to reference production URLs
gulp.task('tmpl', function() {
	var assets = useref.assets();

	return gulp.src('./templates/*.tmpl')
	.pipe(assets)
	.pipe(gulpif('*.js', uglify()))
	.pipe(gulpif('*.css', minifyCss()))
	.pipe(assets.restore())
	.pipe(useref())
	.pipe(replacePath(/static\//g, '/static/'))
	.pipe(gulp.dest(dir));
});

gulp.task('fonts', function () {
	return gulp.src('bower_components/bootstrap/dist/fonts/*')
	.pipe(gulp.dest(dir+'/static/fonts'));
})

gulp.task('html', function () {
	return gulp.src('./templates/html/**/*.html')
	.pipe(gulp.dest(dir+'/html'));
})

gulp.task('js', function() {
	var files = ['./templates/js/**/*.js'];
	
	return gulp.src(files)
	.pipe(concat('coolio.bundle.js'))
	.pipe(gulp.dest(dir+'/js'));
});

gulp.task('libs', function() {
	var files = [
	'bower_components/jquery/dist/jquery.min.js',
	'bower_components/bootstrap/dist/js/bootstrap.min.js',
	'bower_components/knockout/dist/knockout.js'
	];
	
	return gulp.src(files)
	.pipe(concat('coolio.libs.min.js'))
	.pipe(gulp.dest(dir+'/libs'));
});

gulp.task('coolio', function() {
	return gulp.src('./templates/coolio.js')
	.pipe(gulp.dest(dir));
});

// Removes all files from ./dist/
gulp.task('clean', function() {
	del(dir);
});

gulp.task('css', function () {
	var files = [
	'bower_components/bootstrap/dist/css/bootstrap.min.css'
	];
	
	return gulp.src(files)
	.pipe(concat('coolio.min.css'))
	.pipe(gulp.dest(dir+'/static/css'));
});

gulp.task('watch', function() {
	gulp.watch('./templates/coolio.js', ['coolio']);
	gulp.watch('./templates/html/**/*.html', ['html']);
	gulp.watch('./templates/js/**/*.js', ['js']);
	gulp.watch('./templates/*.tmpl', ['tmpl']);
	gulp.watch('./templates/static/**/*', ['tmpl']);
});

gulp.task('default', ['clean', 'tmpl','html', 'css', 'fonts', 'libs', 'js', 'coolio', 'watch']);
