package xcodeproj

const bluePrintScriptContent = `require 'xcodeproj'
require 'json'

require 'xcodeproj'
require 'json'

def workspace?(project_path)
  extname = File.extname(project_path)
  extname == '.xcworkspace'
end
  
def contained_projects(project_path)
  return [project_path] unless workspace?(project_path)

  workspace = Xcodeproj::Workspace.new_from_xcworkspace(project_path)
  workspace_dir = File.dirname(project_path)
  project_paths = []
  workspace.file_references.each do |ref|
    pth = ref.path
    next unless File.extname(pth) == '.xcodeproj'
    next if pth.end_with?('Pods/Pods.xcodeproj')

    project_path = File.expand_path(pth, workspace_dir)
    project_paths << project_path
  end

  project_paths
end

def blueprint_identifier(project_path, scheme_name)
  schemes_by_project = {}

  project_paths = contained_projects(project_path)
  project_paths.each do |path|
    scheme_path = File.join(path, 'xcshareddata', 'xcschemes', scheme_name + '.xcscheme')
    next unless File.exist?(scheme_path)

    scheme = Xcodeproj::XCScheme.new(scheme_path)

    action = scheme.build_action
    next unless action

    entries = action.entries
    next unless entries

    isArchiving = entries.build_for_archiving
    next unless isArchiving

    references = entries.buildable_references
    next unless references

    bluePrintIdentifier = references.target_uuid
    next unless bluePrintIdentifier

    return bluePrintIdentifier
  end

  raise 'build action Blueprint Identifier not found'
end
  
begin
  path = ENV['PROEJECTPATH']
  ret = blueprint_identifier(path, ENV['SCHEME_NAME'])
  result = {
    data: ret
  }
  result_json = JSON.pretty_generate(result).to_s
  puts result_json
rescue => e
  error_message = e.to_s + "\n" + e.backtrace.join("\n")
  result = {
    error: error_message
  }
  result_json = result.to_json.to_s
  puts result_json
  exit(1)
end
`

const targetInfoGemfileContent = `source "https://rubygems.org"
gem "xcodeproj"
gem "json"
`
