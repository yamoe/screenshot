package io.mine.screenshot;

import javax.validation.constraints.NotNull;

import org.hibernate.validator.constraints.URL;

public class ScreenshotFormData {
	@URL @NotNull
	private String url;
	
	public String getUrl() {
		return this.url;
	}
	
	public void setUrl(String url) {
		this.url = url;
	}

}
