package io.mine.screenshot;

import java.awt.Color;
import java.awt.Graphics2D;
import java.awt.image.BufferedImage;
import java.util.Base64;
import java.util.LinkedHashMap;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.IOException;

import javax.imageio.ImageIO;

import org.openqa.selenium.Dimension;
import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebDriverException;
import org.openqa.selenium.OutputType;
import org.openqa.selenium.TakesScreenshot;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;


public class Screenshot {
	
	private static final String formatName = "png";
	private static final String mimetype = "image/png";
	
	private WebDriver driver;
	private BufferedImage image;

	class Size {
		private int devicePixelRatio;
		private int totalWidth;
		private int viewportWidth;
		private int totalHeight;
		private int viewportHeight;

		public Size(JavascriptExecutor js) {
			devicePixelRatio = ((Long)js.executeScript("return window.devicePixelRatio")).intValue();
		    totalWidth = ((Long)js.executeScript("return document.body.parentNode.scrollWidth")).intValue();
		    viewportWidth = ((Long)js.executeScript("return document.body.clientWidth")).intValue();
		    totalHeight = ((Long)js.executeScript("return document.body.parentNode.scrollHeight")).intValue();
		    viewportHeight = ((Long)js.executeScript("return window.innerHeight")).intValue();
		}
		
		public int getDevicePixelRatio() {
			return devicePixelRatio;
		}

		public int getTotalWidth() {
			return totalWidth;
		}

		public int getViewportWidth() {
			return viewportWidth;
		}

		public int getTotalHeight() {
			return totalHeight;
		}

		public int getViewportHeight() {
			return viewportHeight;
		}
	}
	
	public Screenshot() {}
	
	public void capture(String url, Integer width) throws Exception {
		
		ChromeOptions options = new ChromeOptions();
		options.setHeadless(true);
		options.addArguments("hide-scrollbars");
		
		this.driver = new ChromeDriver(options);
		this.driver.get(url);
		
		try {
			int w = (width != null) ? width.intValue() : 0;
			this.image = this.scrollScreenshot(w);
		} catch (Throwable e) {
			e.printStackTrace();
			throw e;
		} finally {
			driver.quit();
		}
	}
	
	private BufferedImage scrollScreenshot(int width) throws InterruptedException, WebDriverException, IOException {
		long scrollDelay = 300;
		return scrollScreenshot(width, scrollDelay);
	}
	
	
	private BufferedImage scrollScreenshot(int width, long scrollDelay) throws InterruptedException, WebDriverException, IOException {
		//reference : https://stackoverflow.com/questions/41721734/taking-screenshot-of-full-page-with-selenium-python-chromedriver

		JavascriptExecutor js = (JavascriptExecutor)this.driver;
		
		Size size = new Size(js);
		
		int captureWidth = (width > 0) ? width: size.getTotalWidth(); 
		this.driver.manage().window().setSize(new Dimension(captureWidth, size.getViewportHeight()));
		size = new Size(js);

		int totalHeight = size.getTotalHeight();
	    int offset = 0;  // height
	    LinkedHashMap<Integer, BufferedImage> slices = new LinkedHashMap<Integer, BufferedImage>();
	    while (offset < totalHeight) {
	
	        if (offset + size.getViewportHeight() > size.getTotalHeight()) {
	            offset = size.getTotalHeight() - size.getViewportHeight();
	        }
	
	        String str = (new StringBuilder())
	        		.append("window.scrollTo(")
	        		.append(0)
	        		.append(", ")
	        		.append(offset)
	        		.append(")")
	        		.toString();
	        js.executeScript(str);
	
	        Thread.sleep(scrollDelay);
	
	        BufferedImage img = ImageIO.read(
	        	new ByteArrayInputStream((
	        		(TakesScreenshot)this.driver).getScreenshotAs(OutputType.BYTES)
	        	)
	        );
	        slices.put(offset, img);
	        offset = offset + size.getViewportHeight();

	    }
	    
	    BufferedImage stitchedImage = new BufferedImage(
	    		size.getTotalWidth() * size.getDevicePixelRatio(),
	    		size.getTotalHeight() * size.getDevicePixelRatio(),
	    		BufferedImage.TYPE_INT_RGB);
	    
	    Graphics2D graphics = (Graphics2D)stitchedImage.getGraphics();
	    graphics.setBackground(Color.WHITE);
	        
	    for (Integer offsetHeight : slices.keySet()) {
	    	BufferedImage img = slices.get(offsetHeight);
	    	graphics.drawImage(img, 0, offsetHeight * size.getDevicePixelRatio(), null);
	    }
	    return stitchedImage;
	}
	
	public void save(String path) throws Exception {
		ImageIO.write(this.image, Screenshot.formatName, new File(path));
	}
	
	public ByteArrayOutputStream stream() throws Exception {
		ByteArrayOutputStream stream = new ByteArrayOutputStream();
		ImageIO.write(this.image, Screenshot.formatName, stream);
		stream.flush();
		return stream;
	}
	
	public String imgTag() throws Exception {
		ByteArrayOutputStream stream = this.stream();
        byte[] bytes = stream.toByteArray();
        return (new StringBuilder())
        		.append("<img src=\"data:image/png;base64,")
        		.append(Base64.getEncoder().encodeToString(bytes))
        		.append("\">")
        		.toString();
        
	}
	
	public String mimetype() {
		return Screenshot.mimetype;
	}
	
/*
    public void showFile() {
        return send_file(self.io(), mimetype=self.mimetype());
    }
    
    publid void downloadFile(Stringi filename) {
        """browser 에서 바로 다운로드"""
        return send_file(self.io(), mimetype=self.mimetype(),
            as_attachment=True, attachment_filename=filename);
    }
*/


}
