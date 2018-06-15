package io.mine.screenshot;

import java.io.PrintWriter;

import javax.servlet.http.HttpServletResponse;
import javax.validation.Valid;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

import com.google.common.base.Strings;


@Controller
@RequestMapping("/")
public class ScreenshotController {
	

	@GetMapping("/")
	public String main(Model model) {
		model.addAttribute("data", new ScreenshotFormData());
		return "main";
	}
	
	@PostMapping("/")
	public String post(RedirectAttributes redirectAttributes, @Valid @ModelAttribute("data") ScreenshotFormData data, BindingResult result) {
		if (result.hasErrors()) {
			return "main";
		}
	    redirectAttributes.addAttribute("url", data.getUrl());
		return "redirect:/screenshot";
	}
	
	@GetMapping("/screenshot")
	public void screenshot(
			@Valid @ModelAttribute ScreenshotParams params,
			BindingResult result,
			HttpServletResponse response) throws Exception {
		
		if (result.hasErrors()) {
            throw new IllegalArgumentException(result.getAllErrors().toString());
		}

		Screenshot screenshot = new Screenshot();
		screenshot.capture(params.getUrl(), params.getWidth());

		if (params.getTag()) {
			// <img> tag
			PrintWriter out = response.getWriter();
			out.println(screenshot.imgTag());
			return;
		}
		
		String filename = params.getFilename();
		if (Strings.isNullOrEmpty(filename)) {
			// view image
			response.setContentType(screenshot.mimetype());
		} else {
			// download image
		    response.setContentType("application/octet-stream");
		    response.setHeader("Content-Disposition", "attachment;filename=" + filename);
		}
	    response.getOutputStream().write(screenshot.stream().toByteArray());
	}
}
