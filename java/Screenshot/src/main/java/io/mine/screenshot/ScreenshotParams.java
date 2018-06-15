package io.mine.screenshot;

import javax.validation.constraints.NotNull;

import org.hibernate.validator.constraints.URL;

public class ScreenshotParams {
	
	@URL @NotNull
	private String url;
	private Integer width;
	private String filename;
	private Boolean tag = false;
	
	public String getUrl() {
		return url;
	}
	public void setUrl(String url) {
		this.url = url;
	}
	public Integer getWidth() {
		return width;
	}
	public void setWidth(Integer width) {
		this.width = width;
	}
	public String getFilename() {
		return filename;
	}
	public void setFilename(String filename) {
		this.filename = filename;
	}
	public Boolean getTag() {
		return tag;
	}
	public void setTag(Boolean tag) {
		this.tag = tag;
	}
	
}
